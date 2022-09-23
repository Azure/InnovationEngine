from parser import Parser 

def main():
    parser = Parser("README.md")

    for item in parser.markdownElements:
        print(item[0])
        print(item[1].subtype)
        print(item[1].value)
    
    

main()
