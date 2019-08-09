def hex2rgb(hexcolor):
    rgb = [255, 255, 255]
    try:
        rgb = [(hexcolor >> 16) & 0xff, (hexcolor >> 8) & 0xff, hexcolor & 0xff]
    except Exception, e:
        pass
    return rgb


def rgb2hex(rgbcolor):
    r, g, b = rgbcolor
    return (r << 16) + (g << 8) + b


if __name__ == '__main__':
    print("rgb2hex((128,128,18))=%s" % rgb2hex((128, 128, 18)))
    print("hex2rgb(8421394)=%s" % hex2rgb(8421394))
