package af

// relative coordinates range from -3 to +3
func RelativeCoordinateToBlockId(dx int16, dy int16, dz int16) (id uint16) {
    id =  0
    id |= uint16(dx + 3)
    id |= uint16(dy + 3) << 3
    id |= uint16(dz + 3) << 6
    return 
}   

func BlockIdToRelativeCoordinate(id uint16) (dx int16, dy int16, dz int16) {
    dx = int16(id & 0x0007 - 3)
    dy = int16((id & 0x0038) >> 3 - 3)
    dz = int16((id & 0x01C0) >> 6 - 3)
    return
}
