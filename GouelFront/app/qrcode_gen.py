import qrcode
from PIL import Image
import base64
from io import BytesIO

ICON_SIZE = 70
QR_SIZE = 128


def convert_transparency_to_white(image):
    if image.mode != "RGBA":
        image = image.convert("RGBA")
    background = Image.new("RGBA", image.size, "white")
    combined = Image.alpha_composite(background, image)
    return combined.convert("RGB")


def generate_qr_code(data, color="black"):
    qr = qrcode.QRCode(
        version=2,
        error_correction=qrcode.constants.ERROR_CORRECT_H,
        box_size=10,
        border=4,
    )
    qr.add_data(data)
    qr.make(fit=True)
    img = qr.make_image(fill_color=color, back_color="white").convert("RGB")
    return img


def convert_image_to_base64(image):
    buffered = BytesIO()
    image.save(buffered, format="PNG")
    return base64.b64encode(buffered.getvalue()).decode()


def qrcode_total(data: str) -> str:
    qr_image = generate_qr_code(data)
    return convert_image_to_base64(qr_image)
