from MarkdownParser import MarkdownParser 
from executor import Executor
import signal
import time
import sys

def main():

    if sys.argv[2].endswith(".md"):
        parser = MarkdownParser(str(sys.argv[2]))
        signal.signal( signal.SIGALRM, timeoutHandler)
        signal.alarm(2)
    
        try:
            parser.parseMarkdown()
            signal.alarm(0)
        except Exception as ex:
            if str(ex) == "Not Executable":
                # exit(1)
                exit(0)
    else:
        print("Currently Innovation Engine can only parse Markdown. The Input file '" + sys.argv[2] + "' is not a markdown file." )
        # Removing exit(1) as it will prematurely terminate the github action if there are multiple files being tested
        # exit(1)

    executor = Executor(parser.markdownElements, str(sys.argv[1]), str(sys.argv[2]))
    executor.runMainLoop()


def timeoutHandler(num, stack):
        print("Document named " + sys.argv[2] + " could not be parsed \n")
        print("Check for errors in the markdown file like non closed codeblocks")
        raise Exception("Not Executable")


main()

