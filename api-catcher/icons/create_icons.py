from PIL import Image, ImageDraw, ImageFont
import os

def create_icon(size, filename):
    # 创建渐变背景
    img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 绘制圆角矩形背景（渐变效果）
    padding = max(1, size // 16)
    
    # 创建渐变色彩
    for y in range(padding, size - padding):
        ratio = (y - padding) / max(1, size - 2 * padding)
        r = int(102 + (118 - 102) * ratio)  # 从#667eea到#764ba2
        g = int(126 + (75 - 126) * ratio)
        b = int(234 + (162 - 234) * ratio)
        draw.line([(padding, y), (size - padding, y)], fill=(r, g, b, 255))
    
    # 绘制白色"API"文字（小尺寸不绘制文字）
    if size >= 48:
        try:
            font_size = max(8, size // 3)
            font = ImageFont.truetype("/System/Library/Fonts/Helvetica.ttc", font_size)
        except:
            font = ImageFont.load_default()
        
        text = "API"
        try:
            bbox = draw.textbbox((0, 0), text, font=font)
            text_width = bbox[2] - bbox[0]
            text_height = bbox[3] - bbox[1]
        except:
            text_width = size // 2
            text_height = size // 3
        
        x = (size - text_width) // 2
        y = (size - text_height) // 2 - max(1, size // 16)
        
        draw.text((x, y), text, fill=(255, 255, 255, 255), font=font)
    
    # 保存
    img.save(filename, 'PNG')
    print(f"Created: {filename} ({size}x{size})")

# 创建三个尺寸的图标
create_icon(16, 'icon16.png')
create_icon(48, 'icon48.png')
create_icon(128, 'icon128.png')

print("All icons created successfully!")
