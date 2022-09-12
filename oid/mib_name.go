package oid

// MibName maps mib name to []uint32
var MibName = map[string][]uint32{
	"internet":     {1, 3, 6, 1},
	"directory":    {1, 3, 6, 1, 1},
	"mgmt":         {1, 3, 6, 1, 2},
	"experimental": {1, 3, 6, 1, 3},
	"private":      {1, 3, 6, 1, 4},
	"security":     {1, 3, 6, 1, 5},
	"snmpV2":       {1, 3, 6, 1, 6},
	"mib2":         {1, 3, 6, 1, 2, 1},
	"enterprises":  {1, 3, 6, 1, 4, 1},
	// System
	"system":          {1, 3, 6, 1, 2, 1},
	"sysDescr":        {1, 3, 6, 1, 2, 1, 1}, // OctetString  read-only
	"sysOid":          {1, 3, 6, 1, 2, 1, 2}, // ObjectID     read-only
	"sysUpTime":       {1, 3, 6, 1, 2, 1, 3}, // TimeTicks    read-only
	"sysContact":      {1, 3, 6, 1, 2, 1, 4}, // OctetString  read-write
	"sysName":         {1, 3, 6, 1, 2, 1, 5}, // OctetString  read-write
	"sysLocation":     {1, 3, 6, 1, 2, 1, 6}, // OctetString  read-write
	"sysServices":     {1, 3, 6, 1, 2, 1, 7}, // Integer32    read-only
	"sysORLastChange": {1, 3, 6, 1, 2, 1, 8}, // TimeTicks    read-only
	"sysORTable":      {1, 3, 6, 1, 2, 1, 9},
	"sysOREntry":      {1, 3, 6, 1, 2, 1, 9, 1},
	"sysORIndex":      {1, 3, 6, 1, 2, 1, 9, 1, 1}, // Integer32    not-accessible
	"sysORID":         {1, 3, 6, 1, 2, 1, 9, 1, 2}, // ObjectID     read-only
	"sysORDescr":      {1, 3, 6, 1, 2, 1, 9, 1, 3}, // OctetString  read-only
	"sysORUpTime":     {1, 3, 6, 1, 2, 1, 9, 1, 4}, // TimeTicks    read-only
	// Interfaces
	"interfaces":    {1, 3, 6, 1, 2, 1, 2},
	"ifNumber":      {1, 3, 6, 1, 2, 1, 2, 1}, // Integer32  read-only
	"ifTable":       {1, 3, 6, 1, 2, 1, 2, 2},
	"ifEntry":       {1, 3, 6, 1, 2, 1, 2, 2, 1},
	"ifIndex":       {1, 3, 6, 1, 2, 1, 2, 2, 1, 1}, // Integer32    read-only
	"ifDescr":       {1, 3, 6, 1, 2, 1, 2, 2, 1, 2}, // OctetString  read-only
	"ifType":        {1, 3, 6, 1, 2, 1, 2, 2, 1, 3}, // Integer32    read-only
	"ifMtu":         {1, 3, 6, 1, 2, 1, 2, 2, 1, 4}, // Integer32    read-only
	"ifSpeed":       {1, 3, 6, 1, 2, 1, 2, 2, 1, 5}, // Unsigned32   read-only
	"ifPhysAddress": {1, 3, 6, 1, 2, 1, 2, 2, 1, 6}, // OctetString  read-only
	"ifAdminStatus": {1, 3, 6, 1, 2, 1, 2, 2, 1, 7}, // Integer32    read-write
	// to be continued ...
}
