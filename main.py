from parser import Parser 
from executor import Executor

def main():
    parser = Parser("README.md")
    parser.parseMarkdown()

    # for item in parser.markdownElements:
    #     print(item[0])
    #     print(item[1].subtype)
    #     print(item[1].value)
    
    executor = Executor(parser.markdownElements)
    executor.runMainLoop()


main()
