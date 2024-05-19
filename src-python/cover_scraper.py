import pypdf

def extract_text_from_cover(pdf_path):
    with open(pdf_path, 'rb') as file:
        pdf_reader = pypdf.PdfReader(file)
        num_pages = len(pdf_reader.pages)
        extracted_text = ''
        
        for page_num in range(num_pages):
            page = pdf_reader.pages[page_num]
            page_text = page.extract_text()
            extracted_text += page_text

    cover_text = ""
    for char in extracted_text:
        if char == '\n':
            cover_text += ' '
        else:
            cover_text += char

    return cover_text
