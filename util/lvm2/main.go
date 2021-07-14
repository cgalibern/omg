// +build linux

package lvm2

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"opensvc.com/opensvc/util/command"
	"opensvc.com/opensvc/util/device"
	"opensvc.com/opensvc/util/file"
	"opensvc.com/opensvc/util/funcopt"
	"opensvc.com/opensvc/util/sizeconv"
)

type (
	LVData struct {
		Report []LVReport `json:"report"`
	}
	LVReport struct {
		LV []LVInfo `json:"lv"`
	}
	LVInfo struct {
		LVName          string `json:"lv_name"`
		VGName          string `json:"vg_name"`
		LVAttr          string `json:"lv_attr"`
		LVSize          string `json:"lv_name"`
		Origin          string `json:"origin"`
		DataPercent     string `json:"data_percent"`
		CopyPercent     string `json:"copy_percent"`
		MetadataPercent string `json:"metadata_percent"`
		MovePV          string `json:"move_pv"`
		ConvertPV       string `json:"convert_pv"`
		MirrorLog       string `json:"mirror_log"`
		Devices         string `json:"devices"`
	}
	driver struct{}
	LV     struct {
		driver
		LVName string
		VGName string
		log    *zerolog.Logger
	}
	LVAttrIndex uint8
	LVAttrs     string
	LVAttr      rune
)

const (
	LVAttrIndexType        LVAttrIndex = 0
	LVAttrIndexPermissions LVAttrIndex = iota
	LVAttrIndexAllocationPolicy
	LVAttrIndexAllocationFixedMinor
	LVAttrIndexState
	LVAttrIndexDeviceOpen
	LVAttrIndexTargetType
	LVAttrIndexZeroDataBlocks
	LVAttrIndexVolumeHealth
	LVAttrIndexSkipActivation
)

const (
	// State attrs field (index 4)
	LVAttrStateActive                               LVAttr = 'a'
	LVAttrStateHistorical                           LVAttr = 'h'
	LVAttrStateSuspended                            LVAttr = 's'
	LVAttrStateInvalidSnapshot                      LVAttr = 'I'
	LVAttrStateSuspendedSnapshot                    LVAttr = 'S'
	LVAttrStateSnapshotMergeFailed                  LVAttr = 'm'
	LVAttrStateSuspendedSnapshotMergeFailed         LVAttr = 'M'
	LVAttrStateMappedDevicePresentWithoutTable      LVAttr = 'd'
	LVAttrStateMappedDevicePresentWithInactiveTable LVAttr = 'i'
	LVAttrStateThinPoolCheckNeeded                  LVAttr = 'c'
	LVAttrStateSuspendedThinPoolCheckNeeded         LVAttr = 'C'
	LVAttrStateUnknown                              LVAttr = 'X'
)

var (
	ErrExist = errors.New("lv does not exist")
)

func (t driver) DriverName() string {
	return "lvm2"
}

func NewLV(vg string, lv string, opts ...funcopt.O) *LV {
	t := LV{
		VGName: vg,
		LVName: lv,
	}
	_ = funcopt.Apply(&t, opts...)
	return &t
}
func WithLogger(log *zerolog.Logger) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*LV)
		t.log = log
		return nil
	})
}

func (t LV) FQN() string {
	return fmt.Sprintf("%s/%s", t.VGName, t.LVName)
}

func (t LV) DevPath() string {
	return fmt.Sprintf("/dev/%s/%s", t.VGName, t.LVName)
}

func (t *LV) Activate() error {
	return t.change([]string{"-ay"})
}

func (t *LV) Deactivate() error {
	return t.change([]string{"-an"})
}

func (t *LV) change(args []string) error {
	fqn := t.FQN()
	cmd := command.New(
		command.WithName("lvchange"),
		command.WithArgs(append(args, fqn)),
		command.WithLogger(t.log),
		command.WithCommandLogLevel(zerolog.InfoLevel),
		command.WithStdoutLogLevel(zerolog.InfoLevel),
		command.WithStderrLogLevel(zerolog.ErrorLevel),
	)
	cmd.Run()
	if cmd.ExitCode() != 0 {
		return fmt.Errorf("%s error %d", cmd, cmd.ExitCode())
	}
	return nil
}

func (t *LV) Show() (*LVInfo, error) {
	data := LVData{}
	fqn := t.FQN()
	cmd := command.New(
		command.WithName("lvs"),
		command.WithVarArgs("--reportformat", "json", fqn),
		command.WithLogger(t.log),
		command.WithCommandLogLevel(zerolog.DebugLevel),
		command.WithStdoutLogLevel(zerolog.DebugLevel),
		command.WithStderrLogLevel(zerolog.DebugLevel),
		command.WithBufferedStdout(),
	)
	if err := cmd.Run(); err != nil {
		if cmd.ExitCode() == 5 {
			return nil, errors.Wrap(ErrExist, fqn)
		}
		return nil, err
	}
	if err := json.Unmarshal(cmd.Stdout(), &data); err != nil {
		return nil, err
	}
	if len(data.Report) == 1 && len(data.Report[0].LV) == 1 {
		return &data.Report[0].LV[0], nil
	}
	return nil, errors.Wrap(ErrExist, fqn)
}

func (t *LV) Attrs() (LVAttrs, error) {
	lvInfo, err := t.Show()
	switch {
	case errors.Is(err, ErrExist):
		return "", nil
	case err != nil:
		return "", err
	default:
		return LVAttrs(lvInfo.LVAttr), nil
	}
}

func (t LVAttrs) Attr(index LVAttrIndex) LVAttr {
	if len(t) < int(index)+1 {
		return ' '
	}
	return LVAttr(t[index])
}

func (t *LV) Exists() (bool, error) {
	_, err := t.Show()
	switch {
	case errors.Is(err, ErrExist):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (t *LV) IsActive() (bool, error) {
	if attrs, err := t.Attrs(); err != nil {
		return false, err
	} else {
		return attrs.Attr(LVAttrIndexState) == LVAttrStateActive, nil
	}
}

func (t *LV) Devices() ([]*device.T, error) {
	l := make([]*device.T, 0)
	data := LVData{}
	fqn := t.FQN()
	cmd := command.New(
		command.WithName("lvs"),
		command.WithVarArgs("-o", "devices", "--reportformat", "json", fqn),
		command.WithLogger(t.log),
		command.WithStdoutLogLevel(zerolog.DebugLevel),
		command.WithStderrLogLevel(zerolog.DebugLevel),
		command.WithBufferedStdout(),
	)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cmd.Stdout(), &data); err != nil {
		return nil, err
	}
	if len(data.Report) == 0 {
		return nil, fmt.Errorf("%s: no report", cmd)
	}
	switch len(data.Report[0].LV) {
	case 0:
		return nil, fmt.Errorf("lv %s not found", fqn)
	case 1:
		// expected
	default:
		return nil, fmt.Errorf("lv %s has multiple matches", fqn)
	}
	for _, s := range strings.Fields(data.Report[0].LV[0].Devices) {
		path := strings.Split(s, "(")[0]
		dev := device.New(path, device.WithLogger(t.log))
		l = append(l, dev)
	}
	return l, nil
}

func (t *LV) Create(size string, args []string) error {
	if i, err := sizeconv.FromSize(size); err == nil {
		// default unit is not "B", explicitely tell
		size = fmt.Sprintf("%dB", i)
	}
	cmd := command.New(
		command.WithName("lvcreate"),
		command.WithArgs(append(args, "--yes", "-L", size, "-n", t.LVName, t.VGName)),
		command.WithLogger(t.log),
		command.WithCommandLogLevel(zerolog.InfoLevel),
		command.WithStdoutLogLevel(zerolog.InfoLevel),
		command.WithStderrLogLevel(zerolog.ErrorLevel),
	)
	cmd.Run()
	if cmd.ExitCode() != 0 {
		return fmt.Errorf("%s error %d", cmd, cmd.ExitCode())
	}
	return nil
}

func (t *LV) Wipe() error {
	path := t.DevPath()
	if !file.Exists(path) {
		t.log.Info().Msgf("skip wipe: %s does not exist", path)
		return nil
	}
	dev := device.New(path, device.WithLogger(t.log))
	return dev.Wipe()
}

func (t *LV) Remove(args []string) error {
	bdev := t.DevPath()
	cmd := command.New(
		command.WithName("lvremove"),
		command.WithArgs(append(args, bdev)),
		command.WithLogger(t.log),
		command.WithCommandLogLevel(zerolog.InfoLevel),
		command.WithStdoutLogLevel(zerolog.InfoLevel),
		command.WithStderrLogLevel(zerolog.ErrorLevel),
	)
	cmd.Run()
	if cmd.ExitCode() != 0 {
		return fmt.Errorf("%s error %d", cmd, cmd.ExitCode())
	}
	return nil
}
