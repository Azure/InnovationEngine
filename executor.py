
# Class which will run the main loop of the program 

from unittest import result
import pexpect
import pexpect.replwrap
import time

PEXPECT_PROMPT = u'[PEXPECT_PROMPT>'
PEXPECT_CONTINUATION_PROMPT = u'[PEXPECT_PROMPT+'

class Executor:
    shell = None
    markdownData = None
    executableCodeList = None

    def __init__(self, markdownData):
        self.markdownData = markdownData
        self.executableCodeList = {"bash", "terraform", 'azurecli-interactive' , 'azurecli'}
        self.shell = self.get_shell()


    def runMainLoop(self):
        beginningHeading = True
        fromCodeBlock = False
        while len(self.markdownData) > 0:

            if (self.markdownData[0][0] == '#'):
                if beginningHeading or fromCodeBlock:
                    print(self.markdownData[0][1].value)
                    self.markdownData.pop(0)
                    beginningHeading = False
                    fromCodeBlock = False
                
                else:
                    beginningHeading = True
                    print("\n\nPress any key to continue... Press b to exit the program \n \n")
                    keyPressed = self.getInstructionKey()
                    if keyPressed == 'b':
                        print("Exiting program on b key press")
                        break

            elif (self.markdownData[0][0] == 'p'):
                print(self.markdownData[0][1].value)
                self.markdownData.pop(0)


            elif (self.markdownData[0][0] == '```'):
                print('```' + self.markdownData[0][1].subtype + '\n' + self.markdownData[0][1].value + '\n```')
                self.executeCode()
                self.markdownData.pop(0)
                fromCodeBlock = True
                
            else:
                self.markdownData.pop(0)
            

    
    def outputMarkdownElement(self):
        pass

    # Goes through markdown elements until it reaches a code block or next heading 
    def outputHeading(self, headingElement):

        pass

    
    def outputCodeblock(self, codeblockElement):
        pass

    def executeCode(self):
        if self.markdownData[0][1].subtype in self.executableCodeList:
            print("\n\nPress any key to execute the above code block... Press b to exit the program \n \n")
            keyPressed = self.getInstructionKey()
            if keyPressed == 'b':
                print("Exiting program on b key press")
                exit()
            self.runCommand()
           
            
        else:
            print("\n\nPress any key to continue... Press b to exit the program \n \n")
            keyPressed = self.getInstructionKey()
            if keyPressed == 'b':
                print("Exiting program on b key press")
                exit()

    def runCommand(self):
        command = self.markdownData[0][1].value

        print("debug", "Execute command: '" + command + "'\n")
        startTime = time.time()
        response = self.shell.run_command(command)
        timeToExecute = time.time() - startTime

        print(response + "\n" + str(timeToExecute))

    
    def getInstructionKey(self):
        """Waits for a single keypress on stdin.
        This is a silly function to call if you need to do it a lot because it has
        to store stdin's current setup, setup stdin for reading single keystrokes
        then read the single keystroke then revert stdin back after reading the
        keystroke.
        Returns the character of the key that was pressed (zero on
        KeyboardInterrupt which can happen when a signal gets handled)
        This method is licensed under cc by-sa 3.0 
        Thanks to mheyman http://stackoverflow.com/questions/983354/how-do-i-make-python-to-wait-for-a-pressed-key\
        """
        import termios, fcntl, sys, os
        fd = sys.stdin.fileno()
        # save old state
        flags_save = fcntl.fcntl(fd, fcntl.F_GETFL)
        attrs_save = termios.tcgetattr(fd)
        # make raw - the way to do this comes from the termios(3) man page.
        attrs = list(attrs_save) # copy the stored version to update
        # iflag
        attrs[0] &= ~(termios.IGNBRK | termios.BRKINT | termios.PARMRK 
                      | termios.ISTRIP | termios.INLCR | termios. IGNCR 
                      | termios.ICRNL | termios.IXON )
        # oflag
        attrs[1] &= ~termios.OPOST
        # cflag
        attrs[2] &= ~(termios.CSIZE | termios. PARENB)
        attrs[2] |= termios.CS8
        # lflag
        attrs[3] &= ~(termios.ECHONL | termios.ECHO | termios.ICANON
                      | termios.ISIG | termios.IEXTEN)
        termios.tcsetattr(fd, termios.TCSANOW, attrs)
        # turn off non-blocking
        fcntl.fcntl(fd, fcntl.F_SETFL, flags_save & ~os.O_NONBLOCK)
        # read a single keystroke
        try:
            ret = sys.stdin.read(1) # returns a single character
        except KeyboardInterrupt:
            ret = 0
        finally:
            # restore old state
            termios.tcsetattr(fd, termios.TCSAFLUSH, attrs_save)
            fcntl.fcntl(fd, fcntl.F_SETFL, flags_save)
        return ret

        
    def get_shell(self):
        """Gets or creates the shell in which to run commands for the
        supplied demo
        """
        if self.shell == None:
            child = pexpect.spawnu('/bin/bash', echo=False, timeout=None)
            ps1 = PEXPECT_PROMPT[:5] + u'\[\]' + PEXPECT_PROMPT[5:]
            ps2 = PEXPECT_CONTINUATION_PROMPT[:5] + u'\[\]' + PEXPECT_CONTINUATION_PROMPT[5:]
            prompt_change = u"PS1='{0}' PS2='{1}' PROMPT_COMMAND=''".format(ps1, ps2)
            shell = pexpect.replwrap.REPLWrapper(child, u'\$', prompt_change)
        return shell