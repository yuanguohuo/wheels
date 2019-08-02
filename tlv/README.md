# Overview

a TLV object has 3 state:

- Open: the length of the TLV's value is unkonwn, the TLV is open-sized. Depending on its tag type (RawTag or NestingTag), you can fill any lenght of either raw bytes or nested TLV objects into it (as long as the underlying buffer doesn't overflow). After that, you can close it;
- Allocated: the only difference between an **Allocated** TLV and an **Open** TLV is that, the space for the **Allocated** TLV's value is pre-allocated; you can fill in raw-bytes or nested TLV objects (depending on its tag type) but the total length you fill in must be exactly the same as the pre-allocated space;
- Closed: the TLV is completed and immutable.

# Examples

## Closed TLV

```
buf := make([]byte, 128)

tag := tlv.RawTag(1)
t, _ := tlv.Build(buf, nil, tag, []byte("ABCD"))

//tag:1426063361 len:4 val:ABCD
fmt.Printf("tag:%d len:%d val:%s\n", t.Tag(), t.Len(), string(t.Val()))
```

## Allocated TLV

```
buf := make([]byte, 128)

rootTag := tlv.NestingTag(10)
root, _ := tlv.Allocate(buf, nil, rootTag, 32)

//total size of son1: TAG_SIZE + LEN_SIZE + 4 = 12
son1Tag := tlv.RawTag(11)
son1, _ := root.OpenChild(son1Tag)
son1.FillRaw([]byte("AABB"))
son1.Close()

//tag:1426063371 len:4 val:AABB
fmt.Printf("tag:%d len:%d val:%s\n", son1.Tag(), son1.Len(), string(son1.Val()))

//total size of son2: TAG_SIZE + LEN_SIZE + 12 = 20
son2Tag := tlv.RawTag(12)
son2, _ := root.BuildChild(son2Tag, []byte("AAAABBBBCCCC"))

//tag:1426063372 len:12 val:AAAABBBBCCCC
fmt.Printf("tag:%d len:%d val:%s\n", son2.Tag(), son2.Len(), string(son2.Val()))

//32 bytes was allocated for root, and we have filled in 32 bytes; so Close()
//will succeed here; otherwise, an error will be returned;
root.Close()
```

## Open TLV

```
buf := make([]byte, 128)

rootTag := tlv.NestingTag(1)
root, err := tlv.Open(buf, nil, rootTag)
if err != nil {
	log.Panicf(err.Error())
}

son1Tag := tlv.NestingTag(2)
son1, err := root.OpenChild(son1Tag)
if err != nil {
	log.Panicf(err.Error())
}
son11Tag := tlv.RawTag(3)
son11, err := son1.AllocateChild(son11Tag, 9)
if err != nil {
	log.Panicf(err.Error())
}
if err := son11.FillRaw([]byte("AAAAAAAAA")); err != nil {
	log.Panicf(err.Error())
}
if err := son11.Close(); err != nil {
	log.Panicf(err.Error())
}

son12Tag := tlv.RawTag(4)
if _, err := son1.BuildChild(son12Tag, []byte("BBBB")); err != nil {
	log.Panicf(err.Error())
}

if err := son1.Close(); err != nil {
	log.Panicf(err.Error())
}
if err := root.Close(); err != nil {
	log.Panicf(err.Error())
}

fmt.Println(root.Val())
//[2 0 0 170 29 0 0 0 3 0 0 85 9 0 0 0 65 65 65 65 65 65 65 65 65 4 0 0 85 4 0 0 0 66 66 66 66]
// --------- -------- -------- ------- -------------------------- -------- ------- -----------
// son1Tag   son1Len  son11Tag son11Len      son11Val             son12Tag son12Len  son12Val
```
