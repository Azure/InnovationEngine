# author: Naman Parikh
# team: Getting Great At Linux (GGAL)
# product: Executable Documentation
# tool: AI azure docs to valid exec docs converter 

import os
import openai
import requests
from bs4 import BeautifulSoup
import re
from github import Github
import time
import yaml
import subprocess

### USER INPUT NEEDED ###
openai.api_type = "azure"
openai.api_version = "2023-05-15"
openai.api_base = os.getenv("OPENAI_API_BASE")  # Your Azure OpenAI resource's endpoint value.
openai.api_key = os.getenv("OPENAI_API_KEY") # Create a Github instance using an access token or username and password
github_access_token = os.getenv("GitHub_Token")
g = Github(login_or_token=github_access_token)
relevant_azure_docs = ['https://raw.githubusercontent.com/MicrosoftDocs/azure-docs/main/articles/virtual-machines/linux/quick-create-cli.md']
### END OF USER INPUT ###

allowed_commands_list = ['azurecli','bash', 'terraform', 'azure-cli-interactive', 'console', 'yaml']

system_prompt_generate_exec_doc = """You are a Microsoft Azure expert. You write tutorial documentation designed to help new users use Azure seamlessly. You are concise in your descriptions and strictly follow markdown style of documentation. Adhere to the following format.

You separate each command into a section marked with markdown headers (##). Each CLI command is embedded in a Markdown codeblock. You will use environment variables for any parameters used in the CLI commands.

You will provide expected outputs for each command. This will also be contained within a codeblock.

INPUT: {An Azure Document in markdown format} 

OUTPUT DOC IN MARKDOWN FORMAT: {REGENERATE THE DOC GIVEN IN INPUT IN MARKDOWN FORMAT WITH THE FOLLOWING CONDITIONS
1. You will add expected output block(s). Each of them would be enclosed in a ``` code block along with the type of the code block given at the start. They should show the expected terminal output as if I ran as the command(s) in the terminal.ONLY SHOW EXPECTED OUTPUT BLOCK(S) FOR COMMANDS THAT WOULD RETURN SOMETHING IN THE TERMINAL, OTHERWISE DO NOT SHOW ANYTHING. THE RESULT NEEDS TO BE IN JSON AND STRICTLY GIVEN WITHIN 1 PAIR OF CURLY BRACKETS. 
2. You will replace all the parameters within the CLI commands given in the ``` code blocks with appropriate environment variables. THE NAMES SHOULD BE UNIQUE AND SHOULD BE SIMILAR TO THE FORM MyParameterName where ParameterName is obtained using the name of the parameter in the command.
3. Add a section at the top of the doc called "Environment Variables" that defines those environment variables from step 2. 
4. Remove the section(s) from the end of the doc that go over removing resources created in the doc
5. Add a section at the bottom of the doc called "FAQs" that gives a list of FAQs that are answered by the doc. The first FAQ should be "What is the command-specific breakdown of permissions needed to implement this doc?" and the answer should be a list of commands in the doc and the permissions needed to run each of those commands.}

YOU ARE COMMITTED TO EXCELLENCE. TAKE A DEEP BREATH AND GO!"""

system_prompt_generate_expected_result = """You are an Azure CLI, Terraform, and Bash expert. Adhere to the following format.

INPUT: {An Azure CLI, Bash, or Terraform command in a Bash shell} 

OUTPUT IN JSON FORMAT: {Generate the expected terminal output as if I ran this command in the terminal. THE RESULT NEEDS TO BE IN JSON AND STRICTLY GIVEN WITHIN 1 PAIR OF CURLY BRACKETS.}

YOU ARE COMMITTED TO EXCELLENCE. TAKE A DEEP BREATH AND GO!"""

system_prompt_generate_expected_result_yesorno = """You are an Azure CLI, Terraform, and Bash expert. Adhere to the following format.

INPUT: {An Azure CLI, Bash, or Terraform command in a Bash shell} 

OUTPUT (yes or no): {Would this command return anything in the terminal (THIS DOES NOT INCLUDE A SIMPLE STATEMENT ALONG THE LINES OF: HERE IS THE OUTPUT)? RESULT SHOULD ONLY BE ONE OF [yes, no]}

YOU ARE COMMITTED TO EXCELLENCE. TAKE A DEEP BREATH AND GO!"""

system_prompt_generate_permissions = """You are an Azure Permissions expert. Adhere to the following format.

INPUT: {An Azure command}

OUTPUT: {ONLY THE NAMES OF PERMISSION(S) ON AZURE REQUIRED TO RUN the command in INPUT, IN AN ARRAY OF STRING(S)}

YOU ARE COMMITTED TO EXCELLENCE. TAKE A DEEP BREATH AND GO!"""

system_prompt_generate_azure_doc_text = """You are a Microsoft Azure technical documentation writing expert. Adhere to the following format.

INPUT: markdown-styled text with subtitle blocks describing the code chunks below them. 

DESIRED OUTPUT: reproduce the same markdown-styled text with descriptions between the subtitle blocks and code chunks. The final result should be similar to a doc that can be published on Azure. 

TAKE A DEEP BREATH AND GO!"""

system_prompt_generate_relevant_links = """You are an Azure documentation expert. Adhere to the following format:

INPUT: {user's description of a scenario for which they want to leverage Azure} 

OUTPUT: {a list of existing relevant Azure docs (ONLY THAT USE CLI) that, when followed in sequence, will fulfill this scenario COMPLETELY.}

TAKE A DEEP BREATH AND GO!"""

system_prompt_generate_faq = """You are an Azure documentation expert who is great at generating FAQs. Adhere to the following format:

INPUT: {An Azure Document in markdown format}

OUTPUT IN DICT FORMAT: {a list of frequently asked questions, including the most common errors and associated fixes that people would encounter when running the commands given in INPUT, along with specific answers to those questions. every answer should be specific, tangible, and detailed. it may or may not leverage existing resources. include relevant link(s), wherever necessary, in every answer to get more information. OUTPUT SHOULD STRICTLY BE IN DICT FORMAT, WITH KEY BEING THE QUESTION AND VALUE BEING THE ANSWER}

TAKE A DEEP BREATH AND GO!"""

def extract_raw_azure_doc(url):
    response = requests.get('https://learn.microsoft.com/en-us/azure/static-web-apps/get-started-cli?tabs=vanilla-javascript')
    soup = BeautifulSoup(response.text, 'html.parser')
    links = soup.find_all('a')
    for link in links:
        github_url = link.get('href') 
        if github_url is not None and 'github' in github_url and 'md' in github_url:
            github_url = github_url.replace('blob/', '').replace('github.com', 'raw.githubusercontent.com')
            return(github_url)
    # print([link.get('href') for link in links])
    # print(soup.prettify())
    
def get_ai_generated_exec_doc(text=None):
    response = openai.ChatCompletion.create(
        engine="gpt-35-turbo-16k",
        model="gpt-3.5-turbo",
        messages=[
            {"role": "system", "content": system_prompt_generate_exec_doc},
            {"role": "user", "content": f"INPUT: {text}\nOUTPUT DOC:"}
    ])
    expected_result = f"Results:\n\n<!-- expected_similarity=0.3 -->\n```json\n{response.choices[0].message.content}\n```" #currently hardcoding the expected similarity to be 0.3 and can change later
    return expected_result

def get_azure_docs_by_keyword(keyword):
    # Get the repository object from the URL
    repo = g.get_repo('MicrosoftDocs/azure-docs')

    # Get the contents of the repository
    contents = repo.get_contents('')

    # Initialize an empty list to store URLs of files with 'cli' in the URL
    relevant_azure_docs = {}

    try:
        while contents:
            file_content = contents.pop(0)
            if file_content.type == "dir":
                contents.extend(repo.get_contents(file_content.path))
            else:
                if keyword.lower() in file_content.download_url and 'github' in file_content.download_url and 'doc' in file_content.download_url and 'md' in file_content.download_url:
                    if file_content.name not in relevant_azure_docs:
                        relevant_azure_docs[file_content.name] = file_content.download_url
                   
    except KeyboardInterrupt:
        print(relevant_azure_docs)

    return relevant_azure_docs

def get_expected_result(command):
    generate_expected_result_yesorno = openai.ChatCompletion.create(
    engine="gpt-35-turbo-16k",
    model="gpt-3.5-turbo",
    messages=[
        {"role": "system", "content": system_prompt_generate_expected_result_yesorno},
        {"role": "user", "content": f"INPUT: {command}\nOUTPUT (yes or no):"}
    ])
    if generate_expected_result_yesorno.choices[0].message.content.lower() in ['yes', 'y', 'ye', 'yup', 'yus']:
        response = openai.ChatCompletion.create(
        engine="gpt-35-turbo-16k",
        model="gpt-3.5-turbo",
        messages=[
            {"role": "system", "content": system_prompt_generate_expected_result},
            {"role": "user", "content": f"INPUT: {command}\nOUTPUT IN CURLY BRACKETS:"}
        ])
        expected_result = f"Results:\n\n<!-- expected_similarity=0.3 -->\n```json\n{response.choices[0].message.content}\n```" #currently hardcoding the expected similarity to be 0.3 and can change later
        return expected_result
    else:
        return None
    
def get_faq(text=None):
    response = openai.ChatCompletion.create(
    engine="gpt-35-turbo-16k",
    model="gpt-3.5-turbo",
    messages=[
        {"role": "system", "content": system_prompt_generate_faq},
        {"role": "user", "content": f"INPUT: {text}\nOUTPUT IN DICT FORMAT:"}
    ])
    return response.choices[0].message.content

def modify_header_permissions_ie_tag(text=None, all_permissions=None):
    pattern = r'---\n(.*?)\n---'
    matches = ''.join(re.findall(pattern, text, re.DOTALL)) #[0].split('\n')
    if 'ms.custom' in matches and 'ms.custom: innovation-engine' not in matches:
        updated_matches = matches.replace('ms.custom:', 'ms.custom: innovation-engine,')
    if 'ms.custom' not in matches:
        updated_matches = f"{matches}\nms.custom: innovation-engine"
    if 'ms.permissions' in matches:
        updated_matches = updated_matches.replace('ms.permissions:', f'ms.permissions: {all_permissions}, ')
    if 'ms.permissions' not in matches:
        updated_matches = f"{updated_matches}\nms.permissions: {all_permissions}"
    text = text.replace(matches, updated_matches)
    return text

def remove_clean_up_resources_section(text):
    try:
        pattern = r'#{1,6} ' + re.escape('Clean up resource') + r'(.*?)(?=\n#{2,6} |\Z)'
        text = re.sub(pattern, '', text, flags=re.MULTILINE|re.DOTALL)
        return text
    except:
        return text

def get_permissions(text):
    permissions_dict = {}    
    commands = re.findall(r'```[\w\W]*?```', text, re.DOTALL)
    for raw_command in commands:
        command_type = [element for element in allowed_commands_list if element in raw_command]
        if len(command_type) > 0:    
            response = openai.ChatCompletion.create(
                engine="gpt-35-turbo-16k",
                model="gpt-3.5-turbo",
                messages=[
                    {"role": "system", "content": system_prompt_generate_permissions},
                    {"role": "user", "content": f"INPUT: {raw_command}\nOUTPUT IN LIST FORMAT:"}
                ]).choices[0].message.content
            response = response.replace('-', '').replace(']', '').replace('[', '').replace('(', '').replace(')', '').replace('{', '').replace('}', '').replace('"', '').replace(',','').replace('`', '').replace("'", "")
            response = re.sub(r'\d+\.', '', response)
            output_list = response.split('\n')
            permissions_list = [item.strip() for item in output_list if 'microsoft' in item.lower()]
            permissions_list = list(set(permissions_list))
            if raw_command not in permissions_dict:
                permissions_dict[raw_command] = permissions_list
            else:
                permissions_dict[raw_command] = f"{permissions_dict[raw_command]}\n{permissions_list}"
    return permissions_dict

def insert_env_var_section(text, env_var_dict=None):
    section_match = re.search(r"^#+\s.*$", text, re.MULTILINE)
    before_you_begin_match = re.search(r"^#+\sBefore you begin*$", text, re.MULTILINE)
    if section_match:
        if before_you_begin_match:
            end_match = re.search(r"^#+\s.*$", text[before_you_begin_match.end():], re.MULTILINE)
            if end_match:
                end_index = before_you_begin_match.end() + end_match.start()
            else:
                end_index = len(text)
        else:
            end_match = re.search(r"^#+\s.*$", text[section_match.end():], re.MULTILINE)
            if end_match:
                end_index = section_match.end() + end_match.start()
            else:
                end_index = len(text)
        if len(env_var_dict) > 0:
            env_vars = ""
            for key in env_var_dict:
                env_vars += f"\nexport {key}={env_var_dict[key]}"
        
            text = text[:end_index] + "\n" + f"""## Define Environment Variables\n\nThe First step in this tutorial is to define environment variables.\n\n```bash\nexport RANDOM_ID="$(openssl rand -hex 3)"{env_vars}\n```\n"""+text[end_index:]
    
    return text

def get_environment_vars(text):
    code_dict = {}
    env_var_dict = {}
    commands = re.findall(r'```[\w\W]*?```', text)
    for raw_command in commands:
        command_type = [element for element in allowed_commands_list if element in raw_command]
        if len(command_type) > 0: 
          pattern = r"--(\S+)\s+(\S+)"
          matches = re.findall(pattern, raw_command)
          if len(matches) > 0:
            for match in matches:
                key, value = match[0], match[1]
                if key in code_dict and value not in code_dict[key]:
                    code_dict[key].append(value)
                elif key in code_dict and value in code_dict[key]:
                    continue
                else:
                    code_dict[key] = [value]
    for key in code_dict:
        for index, value in enumerate(code_dict[key]):
            if len(code_dict[key]) > 1:
              env_key = ''.join([word.capitalize() for word in key.split('-')])  
              env_var_dict[f'My{env_key}{index+1}'] = value
            else:
              env_key = ''.join([word.capitalize() for word in key.split('-')])  
              env_var_dict[f'My{env_key}'] = value
    return env_var_dict

def get_azure_doc_text(url): 
    response = requests.get(url) 
    return(response.text) 

def get_latest_error_log():
    with open('ie.log', 'rb') as f:
        f.seek(-5, os.SEEK_END)
        while f.read(5) != b'time=':
            f.seek(-6, os.SEEK_CUR)
        last_log = f.read().decode()
    return last_log

def create_pr(dirname):  
    os.environ['GITHUB_TOKEN'] = github_access_token
    branch_name = f"aiexecdocs/{dirname}" # os.popen('git rev-parse --abbrev-ref HEAD').read().strip()
    label = "ai-doc"
    subprocess.run(f'git checkout -b {branch_name}', shell=True)
    subprocess.run(f'git remote add aiexecdocs https://{github_access_token}@github.com/Azure/InnovationEngine.git', shell=True)
    subprocess.run(f'git add .', shell=True)
    subprocess.run(f'git commit -m "AI generated doc titled {dirname}"', shell=True)
    subprocess.run(f'git push aiexecdocs {branch_name} -f', shell=True)
    # subprocess.run(['gh', 'pr', 'create', '--title', f"AI generated doc titled {dirname}", '--body',  f"This PR was generated by AI for an existing Azure Doc: {dirname}. Please review and merge.", '--label', label])
    process = subprocess.Popen(['gh', 'pr', 'create', '--title', f"[RFC] AI generated doc titled {dirname}", '--body',  f"This PR was generated by AI for an existing Azure Doc: {dirname}. Please review and merge.", '--label', label, '--head', branch_name], stdin=subprocess.PIPE, stdout=subprocess.PIPE)

    process.communicate(input=b'\n')

# relevant_azure_docs = get_azure_docs_by_keyword('cli')

for azure_doc_url in relevant_azure_docs:
    if azure_doc_url is not None:
        azure_doc_url = extract_raw_azure_doc(azure_doc_url)
        azure_doc_text = get_azure_doc_text(azure_doc_url)
        match = re.search(r'title: (.*)', azure_doc_text)
        if match:
            azure_doc_name = match.group(1).replace(' ', '').replace("'", "").replace('"', '').replace(':', '').title()
            if not os.path.exists(os.path.join('scenarios/ocd/AIDocs', azure_doc_name.replace('.md', ''))):
                os.makedirs(os.path.join('scenarios/ocd/AIDocs', azure_doc_name.replace('.md', '')))

        if '```' in azure_doc_text:
            with open(os.path.join('scenarios/ocd/AIDocs', azure_doc_name.replace('.md', ''), f'original-{azure_doc_name}'), 'w') as f:
                f.write(azure_doc_text)

            env_var_dict = {}
            if 'Define Environment Variables' not in azure_doc_text:
                env_var_dict = get_environment_vars(azure_doc_text)
            
            commands = re.findall(r'```[\w\W]*?```', azure_doc_text)
            for raw_command in commands:
                command_type = [element for element in allowed_commands_list if element in raw_command]
                if len(command_type) > 0:    
                    command_type = command_type[0]
                    if '```yaml' in raw_command:
                        command = raw_command.replace(command_type,'').replace('`', '')
                        lines = command.split('\n')
                        if lines[0].strip() == '':
                            lines = lines[1:]
                        command = '\n'.join(lines)

                        pattern = r'[^\s]*\.yaml\b'
                        target = '```yaml'
                        matches = [match for match in re.finditer(pattern, azure_doc_text[0:azure_doc_text.find(command)+len(command)])]
                        target_position = azure_doc_text[0:azure_doc_text.find(command)+len(command)].find(target)
                        closest_match = max(matches, key=lambda match: abs(match.start() - target_position))
                        yaml_filename = str(azure_doc_text[0:azure_doc_text.find(command)+len(command)][closest_match.start():closest_match.end()]).replace('`', '').replace('"', '').replace('*', '')
                        with open(os.path.join('scenarios/ocd/AIDocs', azure_doc_name.replace('.md', ''), yaml_filename), 'w') as file:
                            file.write(command)
                    else:
                        command = ' '.join(raw_command.replace(command_type,'').replace('\n', '').replace('`', '').split())
                        old_raw_command = raw_command
                        if len(env_var_dict) > 0:
                            for key in env_var_dict:
                                pattern = r'\b' + re.escape(env_var_dict[key]) + r'\b'
                                if len(env_var_dict) > 0 and re.search(pattern, raw_command):# env_var_dict[key] in raw_command:
                                    raw_command = re.sub(pattern, f'${key}', raw_command)# raw_command.replace(env_var_dict[key], f'${key}')
                        expected_result = get_expected_result(command)
                        if expected_result is not None:
                            azure_doc_text = azure_doc_text.replace(old_raw_command, f"{raw_command}\n\n{expected_result}")

            permissions_dict = get_permissions(azure_doc_text)
            formatted_permissions_dict_str = ''
            for command, permissions in permissions_dict.items():
                if len(permissions) > 0:
                    command = " ".join(command.replace('\n', ' ').split())
                    formatted_permissions_dict_str += f"\n  - {command}\n"
                    for permission in permissions:
                        formatted_permissions_dict_str += f"\n      - {permission}"
            
            faq_dict = eval(get_faq(azure_doc_text))
            
            azure_doc_text += f"""\n<details>\n<summary><h2>FAQs</h2></summary>\n\n#### Q. What is the command-specific breakdown of permissions needed to implement this doc? \n\nA. _Format: Commands as they appears in the doc | list of unique permissions needed to run each of those commands_\n\n{formatted_permissions_dict_str}"""

            if len(faq_dict) > 0:
                for question, answer in faq_dict.items():
                    azure_doc_text += f"""\n\n#### Q. {question} \n\nA. {answer}\n"""
                azure_doc_text += f"""\n</details>"""
            else:
                azure_doc_text += f"""\n</details>"""

            if len(env_var_dict) > 0:
                for key in env_var_dict:
                    if 'name' in key.lower() or 'resourcegroup' in key.lower():
                        env_var_dict[key] += '$RANDOM_ID'
                azure_doc_text = insert_env_var_section(azure_doc_text, env_var_dict)

            all_permissions = []
            for value in permissions_dict.values():
                if isinstance(value, list):
                    all_permissions.extend(value)
            all_permissions = list(set(all_permissions))
            all_permissions = ', '.join(all_permissions)

            azure_doc_text = modify_header_permissions_ie_tag(azure_doc_text, all_permissions)
            azure_doc_text = remove_clean_up_resources_section(azure_doc_text)

            with open(os.path.join('scenarios/ocd/AIDocs', azure_doc_name.replace('.md', ''), f'README.md'), 'w') as f:
                f.write(azure_doc_text)

for dirpath, dirnames, filenames in os.walk('scenarios/ocd/AIDocs'):
    try:        
        for filename in filenames:
            if '.md' in filename:
                file_path = os.path.join(dirpath, filename)
                ie_test_command = f'./bin/ie test {file_path}'
                subprocess.run(ie_test_command, shell=True)
                create_pr(f"{dirpath.split('/')[-1]}-{filename}")
    except:
        print(f"Error: {get_latest_error_log()}")
