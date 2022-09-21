#Given a markdown file, this should break out the headings, paragraphs, executable commands etc. 
class Parser:
    def __init__(self, markdown_filepath):
        self.markdown_filepath = markdown_filepath

    def parse_markdown(self):
        print("Called parser function with file path " + self.markdown_filepath)