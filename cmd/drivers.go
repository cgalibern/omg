package cmd

import (
	_ "opensvc.com/opensvc/drivers/poolshm"
	_ "opensvc.com/opensvc/drivers/resappforking"
	_ "opensvc.com/opensvc/drivers/resappsimple"
	_ "opensvc.com/opensvc/drivers/resdiskloop"
	_ "opensvc.com/opensvc/drivers/resdisklv"
	_ "opensvc.com/opensvc/drivers/resdiskraw"
	_ "opensvc.com/opensvc/drivers/resfsdir"
	_ "opensvc.com/opensvc/drivers/resfsflag"
	_ "opensvc.com/opensvc/drivers/resfshost"
	_ "opensvc.com/opensvc/drivers/resiphost"
	_ "opensvc.com/opensvc/drivers/resiproute"
	_ "opensvc.com/opensvc/drivers/resvol"
)
