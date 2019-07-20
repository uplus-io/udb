/*
 * Copyright (c) 2019 uplus.io
 */

package data

import "testing"

func TestWarehouse_Group(t *testing.T) {
	warehouse := NewWarehouse()
	warehouse.AddNode(NewNode("192.168.1.1", 1117, 1.0), 3, 2)
	warehouse.AddNode(NewNode("192.168.1.1", 1118, 1.0), 3, 2)
	warehouse.AddNode(NewNode("192.168.1.1", 1119, 1.0), 3, 2)
	warehouse.AddNode(NewNode("192.168.2.1", 1117, 1.0), 3, 2)
	warehouse.AddNode(NewNode("192.168.2.2", 1117, 1.0), 3, 2)
	warehouse.AddNode(NewNode("192.168.2.3", 1117, 1.0), 3, 2)
	warehouse.AddNode(NewNode("192.168.3.1", 1117, 1.0), 3, 2)
	warehouse.AddNode(NewNode("192.168.3.2", 1117, 1.0), 3, 2)
	warehouse.AddNode(NewNode("192.168.3.3", 1117, 1.0), 3, 2)
	warehouse.AddNode(NewNode("192.9.1.1", 1117, 1.0), 3, 2)
	warehouse.AddNode(NewNode("10.1.1.1", 1117, 1.0), 3, 2)
	warehouse.AddNode(NewNode("20.1.1.1", 1117, 1.0), 3, 2)
	warehouse.Group()

	warehouse.print()
}
