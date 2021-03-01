package strings

import (
    `testing`
)

func TestDeDuplication(t *testing.T) {
    r := DeDuplication("aaabbb")
    if r != "ab" {
        t.Errorf("ascii 失败：%s", r)
    }
    
    r = DeDuplication("我我我你你你")
    if r != "我你" {
        t.Errorf("中文失败: %s", r)
    }
    
    r = DeDuplication("😀😂☸😀😂")
    if r != "😀😂☸" {
        t.Errorf("表情失败: %s", r)
    }
}
