from parser import Parser 
from executor import Executor
import sys

def main():
    parser = Parser(str(sys.argv[2]))
    parser.parseMarkdown()
    # test
    executor = Executor(parser.markdownElements, str(sys.argv[1]))
    executor.runMainLoop()

main()
