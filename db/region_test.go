// SPDX-License-Identifier: MIT

package db

import (
	"testing"

	"github.com/issue9/assert"
)

func TestRegion_IsSupported(t *testing.T) {
	a := assert.New(t)

	obj := &DB{versions: []int{2020, 2019, 2018}}
	obj.region = &Region{Items: []*Region{
		{Supported: 3, Name: "test", db: obj},
	}, db: obj}

	a.True(obj.region.Items[0].IsSupported(2020))
	a.True(obj.region.Items[0].IsSupported(2019))
	a.False(obj.region.Items[0].IsSupported(2018)) // 不支持
	a.False(obj.region.Items[0].IsSupported(2009)) // 不存在于 db
}

func TestRegion_AddItem(t *testing.T) {
	a := assert.New(t)

	obj := &DB{versions: []int{2020, 2019, 2018}}
	obj.region = &Region{Items: []*Region{}, db: obj}

	a.ErrorString(obj.region.AddItem("33", "浙江", 2001), "不支持该年份")

	a.NotError(obj.region.AddItem("44", "广东", 2020))
	a.Equal(obj.region.Items[0].ID, "44").
		NotNil(obj.region.Items[0].db).
		True(obj.region.Items[0].IsSupported(2020)).
		False(obj.region.Items[0].IsSupported(2019))

	a.ErrorString(obj.region.AddItem("44", "广东", 2020), "存在相同")
}

func TestRegion_SetSupported(t *testing.T) {
	a := assert.New(t)

	obj := &DB{versions: []int{2020, 2019, 2018}}
	obj.region = &Region{Items: []*Region{{db: obj}}, db: obj}

	a.NotError(obj.region.AddItem("33", "浙江", 2020))
	a.NotError(obj.region.Items[0].SetSupported(2020))
	a.NotError(obj.region.Items[0].SetSupported(2019))
	a.ErrorString(obj.region.Items[0].SetSupported(2001), "不存在该年份")
}

func TestFindEnd(t *testing.T) {
	a := assert.New(t)

	data := []byte("0123{56}")
	a.Equal(findEnd(data), 7)
}
