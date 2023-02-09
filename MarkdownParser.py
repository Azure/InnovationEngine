#Given a markdown file, this should break out the headings, paragraphs, executable commands etc. 
import re

class MarkdownParser:
  

    def __init__(self, markdownFilepath):
        self.markdownFilepath = markdownFilepath
        self.markdownElements = []
        self.codeBlockType = '```'
        self.headingType = '#'
        self.paragraphType = 'p'
        self.commentType = '<!--'

       
    def parseMarkdown(self):
        self.markdownFile = open(self.markdownFilepath)
        special_characters = ["#", "`", "~" ,"-" ]
        
        char = self.markdownFile.read(1)

        while char:

            if char == '#':
                self.processHeading(char)
                
            elif char == '`':
                if self.checkForCodeBlock(char):
                    subtype, command = self.processCodeSample(char)
                    self.createAndAppendElement(self.codeBlockType, subtype, command)

            elif char == '<':
                if self.checkForComment():
                    self.processComment(char)
                    
            # No specific tag, creating paragraphs and continuing to loop to a special char
            
            else:
                self.processParagraph(char)
                
            char = self.markdownFile.read(1)

        self.markdownFile.close()

    # Iterates through line adding each character to heading string. Also collects the heading type.
    # Creates a markdown element storing the subtype and text value. Appends to all markdown elements.
    def processHeading(self, char):
        text = ""
        headingCount = 0
        while char != '\n':
            if char == '#':
                headingCount += 1
           
            text += char
            char = self.markdownFile.read(1)
        
        text += char # Adds the new line character to the heading
        subtype = "h" + str(headingCount)
        self.createAndAppendElement(self.headingType, subtype, text)

    # Iterates through text until 3 back ticks ``` are found signifying the end of code block
    def processCodeSample(self, char):
        command = ""
        subtype = ""
        endOfCodeBlock = False
        self.markdownFile.read(2)
        # Reading through initial backtick line to see what the subtype is
        while char != '\n':
            char = self.markdownFile.read(1)
            subtype += char

        while not endOfCodeBlock:
            if (char == '`'):
                if self.checkForCodeBlock(char):
                    endOfCodeBlock = True
                    # Read the remaining bash ticks
                    self.markdownFile.read(2)          

            else:
                command += char
                char = self.markdownFile.read(1)

        # The following if addresses bug raised in createAKSFile test script 
        if 'EOF' in command:
            command = self.removeSpacesBeforeEOF(command)
            command = self.removeLeadingWhitespace(command)

        # Should we read a new line here?
        subtype = subtype.strip()
        command = command.strip()
       
        return subtype, command
        
     
    # Iterates until we find a heading or back-tick. If a heading is found it creates paragraph
    # and leaves the function. If a code block or comment is found, it creates a paragraph
    # then rewinds a character and returns to the main loop to process the comment or code block
    def processParagraph(self, char):
        paragraph = ""
        while char != '#' and char != '':
            if char == '`':
                if self.checkForCodeBlock(char):
                    self.createAndAppendElement(self.paragraphType, 'paragraph', paragraph.strip())
                    self.markdownFile.seek(currentPosition)
                    return
                   # self.processCodeSample(char)
                else:
                    paragraph += char           
            elif char == '<':
                if self.checkForComment:
                    self.createAndAppendElement(self.paragraphType, 'paragraph', paragraph.strip())
                    self.markdownFile.seek(currentPosition)
                    return
                    #self.processComment(char)
                else:
                    paragraph += char
            else: 
                paragraph += char
            
            currentPosition = self.markdownFile.tell()
            char = self.markdownFile.read(1)
            
        if len(self.markdownElements) != 0 and "prerequisites" in self.markdownElements[-1][1].value.lower():
            self.createAndAppendElement(self.paragraphType, 'prerequisites', paragraph.strip())

        elif len(self.markdownElements) != 0 and "next steps" in self.markdownElements[-1][1].value.lower():
            self.createAndAppendElement(self.paragraphType, 'next steps', paragraph.strip())

        else:
            self.createAndAppendElement(self.paragraphType, 'paragraph', paragraph.strip())

        self.markdownFile.seek(currentPosition)

    def processComment(self, char):
        endOfComment = False
        comment = "<"
     
        while not endOfComment and char != '':
            currentPosition = self.markdownFile.tell()
            if self.markdownFile.read(1) == '-' and self.markdownFile.read(1) == '-' and self.markdownFile.read(1) == '>':
                comment += char + '-->'
                endOfComment = True
                continue
            elif self.markdownFile.read(1) == '`':
                if self.checkForCodeBlock(char):
                    subtype, command = self.processCodeSample(char)
                    self.createAndAppendElement(self.commentType, subtype, command)
                    comment += command
            else:
                self.markdownFile.seek(currentPosition)
                comment += char

            char = self.markdownFile.read(1)

        if "expected_similarity" in comment:
            results,subtype = self.processResultsBlock()
            similarity = re.findall(r"\d*\.?\d+", comment)[0]
            # outputblock = "```output\n" + results + "\n```"
            # self.createAndAppendElement(self.paragraphType, 'paragraph', outputblock.strip())
            # Loops through elements starting at the end looking for the most recent
            # Codeblock markdown item and adds the results 
            for i, markdownElement in reversed(list(enumerate(self.markdownElements))):
                if markdownElement[0] == self.codeBlockType:
                    self.markdownElements[i][1].results = results
                    self.markdownElements[i][1].similarity = similarity
            
                    break
            self.createAndAppendElement(self.codeBlockType, subtype, results)
        else:
            self.createAndAppendElement(self.commentType, None, comment)

   


    def createAndAppendElement(self, type, subtype, value):
        element = MarkdownElement(subtype, value)
        self.markdownElements.append((type, element))

    # Helper function ran when we hit a backtick. It checks the next two characters to see
    # if they are also backticks. If so, we enter a code block and return true
    def checkForCodeBlock(self, char):
        currentPosition = self.markdownFile.tell()
        if self.markdownFile.read(1) == '`' and self.markdownFile.read(1) == '`':
            self.markdownFile.seek(currentPosition)
            return True
        else:
            self.markdownFile.seek(currentPosition)
            return False
    
    def checkForComment(self):
        currentPosition = self.markdownFile.tell()
        if self.markdownFile.read(1) == '!' and self.markdownFile.read(1) == '-' and self.markdownFile.read(1) == '-':
            self.markdownFile.seek(currentPosition)
            return True
        else:
            self.markdownFile.seek(currentPosition)
            return False

    def processResultsBlock(self):
        results = ""
        subtype = ""
        endOfCodeBlock = False

        char = self.markdownFile.read(1)
        
        while char != '`':
            char = self.markdownFile.read(1)
        
        self.markdownFile.read(2)
        while char != '\n':
            char = self.markdownFile.read(1)
            subtype += char

        while not endOfCodeBlock:
            if (char == '`'):
                if self.checkForCodeBlock(char):
                    endOfCodeBlock = True
                    # Read the remaining bash ticks
                    self.markdownFile.read(2)
                
            else:
                results += char
            # Read all 3 back ticks
            char = self.markdownFile.read(1)
        
        return results.strip(), subtype.strip()
    # If we want to add process for bold and italicized text
    def processBoldText(self,char):
        pass

    # If want to process dashes for hidden titles etc. 
    def processDash(self, char):
        pass

    def removeSpacesBeforeEOF(self, command):
        numSpacesToRemove = 0
        for char in reversed(command):
            if char == '\n':
                break
            else:
                numSpacesToRemove += 1
        command = command[:-numSpacesToRemove]
        command += 'EOF'

        return command

    def removeLeadingWhitespace(self, command):
        lines = command.split("\n")
        numberOfSpaces = len(lines[1]) - len(lines[1].lstrip())
        output = []
        for line in lines:
            line = line[numberOfSpaces:] if numberOfSpaces <= len(line) else ""
            output.append(line)
        return "\n".join(output)

class MarkdownElement:

    def __init__(self, subtype, value):
        self.subtype = subtype
        self.value = value
        self.results = None
        self.similarity = 1.0

