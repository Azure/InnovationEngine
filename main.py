from parser import Parser 
from executor import Executor
import sys

def main():
    parser = Parser(str(sys.argv[1]))
    parser.parseMarkdown()

    # for item in parser.markdownElements:
    #     print(item[0])
    #     print(item[1].subtype)
    #     print(item[1].value)
    
    executor = Executor(parser.markdownElements)
    executor.runMainLoop()


main()
