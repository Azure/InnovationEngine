
# Class which will run the main loop of the program 

from unittest import result
import pexpect
import pexpect.replwrap
import time
from fuzzywuzzy import fuzz
from fuzzywuzzy import process
import re
from os.path import exists

PEXPECT_PROMPT = u'[PEXPECT_PROMPT>'
PEXPECT_CONTINUATION_PROMPT = u'[PEXPECT_PROMPT+'

class Executor:
    shell = None
    markdownData = None
    executableCodeList = None

    def __init__(self, markdownData, modeOfOperation, fileName):
        self.markdownData = markdownData
        self.fileName = fileName
        self.executableCodeList = {"bash", "terraform", 'azurecli-interactive' , 'azurecli'}
        self.shell = self.get_shell()
        self.readEnvVariables()
        self.modeOfOperation = modeOfOperation
        self.numberOfTestsPassed = 0
        self.totalNumberOfTests = 0
        self.failedTests = []

    # Fairly straight forward main loop. While markdownData is not empty
    # Checks type for heading, code block, or paragrpah. 
    # If Heading it outputs the heading, pops the item and prompts input from user
    # If paragraph it outputs paragraph and pops item from list and continues with no pause
    # If Code block, it calls ExecuteCode helper function to print and execute the code block

    def runMainLoop(self):
        if self.modeOfOperation == "interactive":
            self.runMainLoopInteractive()
        elif self.modeOfOperation == "test":
            self.runMainLoopTest()
        elif self.modeOfOperation == 'execute':
            self.runMainLoopExecute()
        else:
            self.runMainLoopInteractive()


    # This function loops through the markdown elements in an interactive manner. It pauses and 
    # requests input from the user to continue at every heading and code block
    def runMainLoopInteractive(self):
        beginningHeading = True
        fromCodeBlock = False

        for markdownItem in self.markdownData:
            if markdownItem[0] == '#':
                if beginningHeading or fromCodeBlock:
                    print(markdownItem[1].value)
                    beginningHeading = False
                    fromCodeBlock = False
                else:
                    beginningHeading = True
                    self.askForInput("Press any key to continue...")
    
            elif markdownItem[0] == 'p' and markdownItem[1].subtype == 'prerequisites':
                print(markdownItem[1].value)
                self.askForInput("Press any key to proceed and execute the prerequisites...")
                self.executePrerequisites(markdownItem)
                beginningHeading = True

            elif markdownItem[0] == 'p':
                print(markdownItem[1].value)
                beginningHeading = True
            
            elif markdownItem[0] == '```':
                print('\n```' + markdownItem[1].subtype + '\n' + markdownItem[1].value + '\n```')
                self.executeCode(markdownItem)
                fromCodeBlock = True

    # This function runs through and only looks at code blocks. It executes them and then
    # Looks at the output. It will automatically return exit code 1 if a test fails.
    # Used in GitHub Actions and automated testing scenarios
    def runMainLoopTest(self):
        for markdownItem in self.markdownData:
            if markdownItem[0] == '```':
                print('\n```' + markdownItem[1].subtype + '\n' + markdownItem[1].value + '\n```')
                if markdownItem[1].subtype in self.executableCodeList:
                    self.runCommand(markdownItem)
            elif markdownItem[0] == "<!--" and markdownItem[1].subtype == 'variables':
                print('\nSetting Environment Variables - \n' + markdownItem[1].value + '\n')
                self.runCommand(markdownItem)
            elif markdownItem[0] == 'p' and markdownItem[1].subtype == 'prerequisites':
                self.executePrerequisites(markdownItem)

            elif markdownItem[0] == 'p' and markdownItem[1].subtype == 'next steps':
                self.executeNextSteps(markdownItem)
            
        
        print("\n{} of {} Tests Passed!".format(str(self.numberOfTestsPassed), str(self.totalNumberOfTests)))
        if self.numberOfTestsPassed < self.totalNumberOfTests:
            print("---------FAILED CODE BLOCKS-------------- \n\n")
            for failedTest in self.failedTests:
                print(failedTest[0][1].value + '\n')
                print(failedTest[1])
               
            exit(1)

    # This function runs through and executes not pausing for any input or failing.
    # The primary intention of this is for executing pre requisites 
    def runMainLoopExecute(self):
        for markdownItem in self.markdownData:
            print(markdownItem[1].value)
            if markdownItem[0] == '```':
                print('\n```' + markdownItem[1].subtype + '\n' + markdownItem[1].value + '\n```')
                if markdownItem[1].subtype in self.executableCodeList:
                    self.runCommand(markdownItem)

    # Checks to see if code block is executable i.e, bash, terraform, azurecli-interactive, azurecli
    # If it is it will wait for input and call run command which passes the command to the repl
    def executeCode(self, markdownItem):
        if markdownItem[1].subtype in self.executableCodeList:
            self.askForInput("Press any key to execute the above code block...")
            print("Executing Code...")
            self.runCommand(markdownItem)
            
        else:
            self.askForInput("Press any key to continue...")

    # Function takes a command and uses the shell which was instantiated at run time using the 
    # Local shell information to execute that command. If the user is logged into az cli on 
    # The authentication will carry over to this environment as well 
    def runCommand(self, markdownItem):
        command = markdownItem[1].value
        expectedResult = markdownItem[1].results
        expectedSimilarity = markdownItem[1].similarity

        #print("debug", "Execute command: '" + command + "'\n")
        startTime = time.time()
        try:
            # Setting a 20 minute timeout...Need a better way to discover broken commands
            response = self.shell.run_command(command, 1200).strip()
        except ValueError as ve:
            print("Continuation prompt required for command " + command)
            print(ex)
            response = command + " failed to run"
        except Exception as ex:
            print("command timed out")
            print(ex)
            response =  command + " failed to run"
            

        timeToExecute = time.time() - startTime
        print("\n" + response + "\n" + "Time to Execute - " + str(timeToExecute))

        if expectedResult is not None:
            print("Expected Results - " + expectedResult)
            self.testResponse(response, expectedResult, expectedSimilarity, markdownItem)
            

    def testResponse(self, response, expectedResult, expectedSimilarity, markdownItem):
        # Todo... try to implement more than just fuzzy matching. Can we look and see if the command returned 
        # A warning or an error? Problem I am having is calls can return every type of response... I could 
        # Hard code something for Azure responses, but it wouldn't be extendible
        #print("\n```output\n" + expectedResult + "\n```")

        if self.modeOfOperation == "interactive":
            actualSimilarity = fuzz.ratio(response, expectedResult) / 100

            if actualSimilarity < float(expectedSimilarity):
                print("The output is NOT correct. The remainder of the document may not function properly")
                print("The Actual similarity was {} \n The expected similarity was {}".format(str(actualSimilarity), expectedSimilarity))


            self.askForInput("Press any key to continue...")

        elif self.modeOfOperation == "test":
            self.totalNumberOfTests += 1
            actualSimilarity = fuzz.ratio(response, expectedResult) / 100

            if actualSimilarity < float(expectedSimilarity):
                errorOutput = "Test Failed \n\nThe expected result was - \n" + expectedResult
                errorOutput += "\nthe actual result - \n" + response
                errorOutput += "\nThe Actual similarity was {} \nThe expected similarity was {} \n".format(str(actualSimilarity), expectedSimilarity)
                print(errorOutput)
 
                self.failedTests.append((markdownItem, errorOutput))
            else:
                self.numberOfTestsPassed += 1
            

        # todo: Create testing mode for execute which simply lets the user know if a test fails which one it is
    


           
    def executePrerequisites(self, markdownItem):
        results = re.findall(r'\]\(([^)]+)\)', markdownItem[1].value)

        for markdownFilepath in results:
            if exists(markdownFilepath):
                command = 'python3 main.py execute ' + markdownFilepath
                print("this is the command to execute \n" + command)
                response = self.shell.run_command(command).strip()
                print(response)
            else:
                print("Could not find file named " + markdownFilepath + " Please locate and run this prerequisite manually")

    def executeNextSteps(self, markdownItem):
        print("Found Next Steps....")
        pass
    
    def askForInput(self, inputPrompt):
        print("\n\n" + inputPrompt + " Press b to exit the program \n \n")
        keyPressed = self.getInstructionKey()
        if keyPressed == 'b':
            print("Exiting program on b key press")
            exit()

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

    # Function looks for file named 
    def readEnvVariables(self):
        if exists(self.fileName[:-3] + '.ini'):
            envFile = open(self.fileName[:-3] + '.ini')
            lines = envFile.readlines()

            for line in lines:
                variableName = line.split()[0]
                value = line.split()[2]
                command = variableName + '=' + value
                self.shell.run_command(command).strip()


    def get_shell(self):
        """Creates the shell in which to run commands for the 
            innovation engine 
        """
        if self.shell == None:
            child = pexpect.spawnu('/bin/bash', echo=False, timeout=None)
            ps1 = PEXPECT_PROMPT[:5] + u'\[\]' + PEXPECT_PROMPT[5:]
            ps2 = PEXPECT_CONTINUATION_PROMPT[:5] + u'\[\]' + PEXPECT_CONTINUATION_PROMPT[5:]
            prompt_change = u"PS1='{0}' PS2='{1}' PROMPT_COMMAND=''".format(ps1, ps2)
            shell = pexpect.replwrap.REPLWrapper(child, u'\$', prompt_change)
        return shell