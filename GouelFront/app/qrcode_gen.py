import qrcode
from PIL import Image
import base64
from io import BytesIO
import cairosvg

ICON_SIZE = 70
QR_SIZE = 128


def svg_to_png(svg_path, output_path):
    cairosvg.svg2png(
        url=svg_path,
        write_to=output_path,
        output_height=ICON_SIZE,
        output_width=ICON_SIZE,
    )


def convert_transparency_to_white(image):
    if image.mode != "RGBA":
        image = image.convert("RGBA")
    background = Image.new("RGBA", image.size, "white")
    combined = Image.alpha_composite(background, image)
    return combined.convert("RGB")


def generate_qr_code_with_logo(
    data, color="black", logo_path="./app/static/img/icon.svg"
):
    qr = qrcode.QRCode(
        version=2,
        error_correction=qrcode.constants.ERROR_CORRECT_H,
        box_size=10,
        border=4,
    )
    qr.add_data(data)
    qr.make(fit=True)
    img = qr.make_image(fill_color=color, back_color="white").convert("RGB")

    if logo_path.endswith(".svg"):
        temp_logo_path = "./app/static/img/temp_logo.png"
        svg_to_png(logo_path, temp_logo_path)
        logo = Image.open(temp_logo_path)
    else:
        logo = Image.open(logo_path)

    logo = convert_transparency_to_white(logo)
    logo_size = ICON_SIZE
    logo.thumbnail((logo_size, logo_size))
    logo_position = ((img.size[0] - logo_size) // 2, (img.size[1] - logo_size) // 2)
    img.paste(logo, logo_position)

    return img


def convert_image_to_base64(image):
    buffered = BytesIO()
    image.save(buffered, format="PNG")
    return base64.b64encode(buffered.getvalue()).decode()


def qrcode_total(data: str) -> str:
    qr_image = generate_qr_code_with_logo(data)
    return convert_image_to_base64(qr_image)
