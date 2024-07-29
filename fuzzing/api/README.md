# Medusa Fuzzing API

## Description
This API is served when the fuzzing process starts and is terminated automatically when fuzzing ends.

## Configurable Parameters
Configuration parameters exist for the API and can be provided in the config file under the new `apiConfig` object.

- **enabled**: Whether the API should be enabled.
    - **Default**: False


- **port**: The port where the API should run on. If the provided port is unavailable, the API will be served on the next available port in increments of 1.
  - **Default**: 8080


- **wsUpdateInterval**: The interval with which the API will send updates via websocket connections.
    - **Default**: False

## Routes

### Main Routes

- **GET /env**: Returns the environment information. This endpoint returns the current environment information as a JSON response. The client can make a GET request to this endpoint to retrieve the environment information.

- **GET /fuzzing**: Returns the fuzzing information. This endpoint returns the current fuzzing information as a JSON response. The client can make a GET request to this endpoint to retrieve the fuzzing information.

- **GET /logs**: Returns the logs. This endpoint returns the logs as a JSON response. The client can make a GET request to this endpoint to retrieve the logs.

- **GET /coverage**: Returns the coverage information. This endpoint returns the current coverage information as a JSON response. The client can make a GET request to this endpoint to retrieve the coverage information.

- **GET /corpus**: Returns the corpus information. This endpoint returns the current corpus information as a JSON response. The client can make a GET request to this endpoint to retrieve the corpus information.


### Websocket Routes

- **GET /ws/env**: Handles WebSocket connections for the environment. The WebSocket connection is used to stream real-time updates of the environment information. The client can connect to this endpoint to receive updates whenever the environment changes.


- **GET /ws/fuzzing**: Handles WebSocket connections for fuzzing. The WebSocket connection is used to stream real-time updates of the fuzzing information. The client can connect to this endpoint to receive updates whenever the fuzzing status changes.


- **GET /ws/logs**: Handles WebSocket connections for logs. The WebSocket connection is used to stream real-time updates of the logs. The client can connect to this endpoint to receive updates whenever new logs are generated.


- **GET /ws/coverage**: Handles WebSocket connections for coverage. The WebSocket connection is used to stream real-time updates of the coverage information. The client can connect to this endpoint to receive updates whenever the coverage changes.


- **GET /ws/corpus**: Handles WebSocket connections for the corpus. The WebSocket connection is used to stream real-time updates of the corpus. The client can connect to this endpoint to receive updates whenever the corpus changes.


- **GET /ws**: Handles WebSocket connections. This endpoint is a catch-all for all other WebSocket routes. The client can connect to this endpoint to receive updates for all the available WebSocket routes.

### Example return data

- Env info
```json
{
    "config": {
        "fuzzing": {
            "workers": 10,
            "workerResetLimit": 50,
            "timeout": 0,
            "testLimit": 0,
            "shrinkLimit": 5000,
            "callSequenceLength": 100,
            "corpusDirectory": "corpus",
            "coverageEnabled": true,
            "targetContracts": [
                "InnerDeploymentFactory"
            ],
            "targetContractsBalances": [],
            "constructorArgs": {},
            "deployerAddress": "0x30000",
            "senderAddresses": [
                "0x10000",
                "0x20000",
                "0x30000"
            ],
            "blockNumberDelayMax": 60480,
            "blockTimestampDelayMax": 604800,
            "blockGasLimit": 125000000,
            "transactionGasLimit": 12500000,
            "testing": {
                "stopOnFailedTest": true,
                "stopOnFailedContractMatching": false,
                "stopOnNoTests": true,
                "testAllContracts": true,
                "traceAll": false,
                "assertionTesting": {
                    "enabled": true,
                    "testViewMethods": false,
                    "panicCodeConfig": {
                        "failOnCompilerInsertedPanic": false,
                        "failOnAssertion": true,
                        "failOnArithmeticUnderflow": false,
                        "failOnDivideByZero": false,
                        "failOnEnumTypeConversionOutOfBounds": false,
                        "failOnIncorrectStorageAccess": false,
                        "failOnPopEmptyArray": false,
                        "failOnOutOfBoundsArrayAccess": false,
                        "failOnAllocateTooMuchMemory": false,
                        "failOnCallUninitializedVariable": false
                    }
                },
                "propertyTesting": {
                    "enabled": true,
                    "testPrefixes": [
                        "property_"
                    ]
                },
                "optimizationTesting": {
                    "enabled": true,
                    "testPrefixes": [
                        "optimize_"
                    ]
                }
            },
            "chainConfig": {
                "codeSizeCheckDisabled": true,
                "cheatCodes": {
                    "cheatCodesEnabled": true,
                    "enableFFI": false
                }
            }
        },
        "compilation": {
            "platform": "crytic-compile",
            "platformConfig": {
                "target": "nested_deployments.sol",
                "solcVersion": "",
                "exportDirectory": "",
                "args": []
            }
        },
        "logging": {
            "level": "info",
            "logDirectory": "",
            "noColor": false
        },
        "apiConfig": {
            "enabled": true,
            "port": 8080,
            "wsUpdateInterval": 0
        }
    },
    "solcVersion": "0.8.17",
    "system": [
        "GJS_DEBUG_TOPICS=JS ERROR;JS LOG",
        "LANGUAGE=en_NG:en",
        "USER=aaliyah",
        "XDG_SEAT=seat0",
        "XDG_SESSION_TYPE=x11",
        "SHLVL=1",
        "HOME=/home/aaliyah",
        "DESKTOP_SESSION=cinnamon",
        "GIO_LAUNCHED_DESKTOP_FILE=/usr/share/applications/code.desktop",
        "GTK_MODULES=gail:atk-bridge",
        "XDG_SEAT_PATH=/org/freedesktop/DisplayManager/Seat0",
        "DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus",
        "CINNAMON_VERSION=6.0.4",
        "GIO_LAUNCHED_DESKTOP_FILE_PID=1340770",
        "QT_QPA_PLATFORMTHEME=qt5ct",
        "LOGNAME=aaliyah",
        "XDG_SESSION_CLASS=user",
        "XDG_SESSION_ID=c2",
        "GNOME_DESKTOP_SESSION_ID=this-is-deprecated",
        "PATH=/home/aaliyah/.local/bin:/home/aaliyah/.cache/cloud-code/m2c/bin:/home/aaliyah/.nvm/versions/node/v21.7.2/bin:/home/aaliyah/.local/bin:/home/aaliyah/.cargo/bin:/home/aaliyah/.local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin:/home/aaliyah/Desktop/BRAINY_CODELAB/medusa-playground:/home/aaliyah/.avm/bin:/home/aaliyah/Desktop/BRAINY_CODELAB/medusa-playground:/home/aaliyah/.avm/bin",
        "GDM_LANG=en_US",
        "GTK3_MODULES=xapp-gtk3-module",
        "SESSION_MANAGER=local/aaliyah-Latitude-E5570:@/tmp/.ICE-unix/1442,unix/aaliyah-Latitude-E5570:/tmp/.ICE-unix/1442",
        "XDG_SESSION_PATH=/org/freedesktop/DisplayManager/Session0",
        "XDG_RUNTIME_DIR=/run/user/1000",
        "DISPLAY=:0",
        "LANG=en_US.UTF-8",
        "XDG_CURRENT_DESKTOP=X-Cinnamon",
        "XDG_SESSION_DESKTOP=cinnamon",
        "XAUTHORITY=/home/aaliyah/.Xauthority",
        "XDG_GREETER_DATA_DIR=/var/lib/lightdm-data/aaliyah",
        "SSH_AUTH_SOCK=/run/user/1000/keyring/ssh",
        "SHELL=/usr/bin/zsh",
        "QT_ACCESSIBILITY=1",
        "GDMSESSION=cinnamon",
        "GPG_AGENT_INFO=/run/user/1000/gnupg/S.gpg-agent:0:1",
        "GJS_DEBUG_OUTPUT=stderr",
        "XDG_VTNR=7",
        "PWD=/home/aaliyah/Desktop/BRAINY_CODELAB/testing",
        "XDG_CONFIG_DIRS=/etc/xdg/xdg-cinnamon:/etc/xdg",
        "XDG_DATA_DIRS=/usr/share/cinnamon:/usr/share/gnome:/home/aaliyah/.local/share/flatpak/exports/share:/var/lib/flatpak/exports/share:/usr/local/share:/usr/share:/var/lib/snapd/desktop",
        "CHROME_DESKTOP=code-url-handler.desktop",
        "ORIGINAL_XDG_CURRENT_DESKTOP=X-Cinnamon",
        "GDK_BACKEND=x11",
        "OLDPWD=/home/aaliyah/Desktop/BRAINY_CODELAB/testing",
        "ZSH=/home/aaliyah/.oh-my-zsh",
        "PAGER=less",
        "LESS=-R",
        "LSCOLORS=Gxfxcxdxbxegedabagacad",
        "LS_COLORS=rs=0:di=01;34:ln=01;36:mh=00:pi=40;33:so=01;35:do=01;35:bd=40;33;01:cd=40;33;01:or=40;31;01:mi=00:su=37;41:sg=30;43:ca=30;41:tw=30;42:ow=34;42:st=37;44:ex=01;32:*.tar=01;31:*.tgz=01;31:*.arc=01;31:*.arj=01;31:*.taz=01;31:*.lha=01;31:*.lz4=01;31:*.lzh=01;31:*.lzma=01;31:*.tlz=01;31:*.txz=01;31:*.tzo=01;31:*.t7z=01;31:*.zip=01;31:*.z=01;31:*.dz=01;31:*.gz=01;31:*.lrz=01;31:*.lz=01;31:*.lzo=01;31:*.xz=01;31:*.zst=01;31:*.tzst=01;31:*.bz2=01;31:*.bz=01;31:*.tbz=01;31:*.tbz2=01;31:*.tz=01;31:*.deb=01;31:*.rpm=01;31:*.jar=01;31:*.war=01;31:*.ear=01;31:*.sar=01;31:*.rar=01;31:*.alz=01;31:*.ace=01;31:*.zoo=01;31:*.cpio=01;31:*.7z=01;31:*.rz=01;31:*.cab=01;31:*.wim=01;31:*.swm=01;31:*.dwm=01;31:*.esd=01;31:*.jpg=01;35:*.jpeg=01;35:*.mjpg=01;35:*.mjpeg=01;35:*.gif=01;35:*.bmp=01;35:*.pbm=01;35:*.pgm=01;35:*.ppm=01;35:*.tga=01;35:*.xbm=01;35:*.xpm=01;35:*.tif=01;35:*.tiff=01;35:*.png=01;35:*.svg=01;35:*.svgz=01;35:*.mng=01;35:*.pcx=01;35:*.mov=01;35:*.mpg=01;35:*.mpeg=01;35:*.m2v=01;35:*.mkv=01;35:*.webm=01;35:*.webp=01;35:*.ogm=01;35:*.mp4=01;35:*.m4v=01;35:*.mp4v=01;35:*.vob=01;35:*.qt=01;35:*.nuv=01;35:*.wmv=01;35:*.asf=01;35:*.rm=01;35:*.rmvb=01;35:*.flc=01;35:*.avi=01;35:*.fli=01;35:*.flv=01;35:*.gl=01;35:*.dl=01;35:*.xcf=01;35:*.xwd=01;35:*.yuv=01;35:*.cgm=01;35:*.emf=01;35:*.ogv=01;35:*.ogx=01;35:*.aac=00;36:*.au=00;36:*.flac=00;36:*.m4a=00;36:*.mid=00;36:*.midi=00;36:*.mka=00;36:*.mp3=00;36:*.mpc=00;36:*.ogg=00;36:*.ra=00;36:*.wav=00;36:*.oga=00;36:*.opus=00;36:*.spx=00;36:*.xspf=00;36:",
        "NVM_DIR=/home/aaliyah/.nvm",
        "NVM_CD_FLAGS=-q",
        "NVM_BIN=/home/aaliyah/.nvm/versions/node/v21.7.2/bin",
        "NVM_INC=/home/aaliyah/.nvm/versions/node/v21.7.2/include/node",
        "_=/home/aaliyah/Desktop/BRAINY_CODELAB/testing/./medusa",
        "TERM_PROGRAM=vscode",
        "TERM_PROGRAM_VERSION=1.89.1",
        "COLORTERM=truecolor",
        "GIT_ASKPASS=/usr/share/code/resources/app/extensions/git/dist/askpass.sh",
        "VSCODE_GIT_ASKPASS_NODE=/usr/share/code/code",
        "VSCODE_GIT_ASKPASS_EXTRA_ARGS=",
        "VSCODE_GIT_ASKPASS_MAIN=/usr/share/code/resources/app/extensions/git/dist/askpass-main.js",
        "VSCODE_GIT_IPC_HANDLE=/run/user/1000/vscode-git-59f4ae660e.sock",
        "VSCODE_INJECTION=1",
        "ZDOTDIR=/home/aaliyah",
        "USER_ZDOTDIR=/home/aaliyah",
        "TERM=xterm-256color"
    ]
}
```
  
- Fuzzing info
```json
{
    "metrics": {},
    "testCases": [
        {
            "ID": "PROPERTY-InnerInnerDeployment-property-inner-inner-deployment()",
            "LogMessage": {},
            "Message": "[RUNNING] Property Test: InnerInnerDeployment.property_inner_inner_deployment()",
            "Name": "Property Test: InnerInnerDeployment.property_inner_inner_deployment()",
            "Status": "RUNNING"
        },
        {
            "ID": "PROPERTY-InnerDeployment-property-inner-deployment()",
            "LogMessage": {},
            "Message": "[RUNNING] Property Test: InnerDeployment.property_inner_deployment()",
            "Name": "Property Test: InnerDeployment.property_inner_deployment()",
            "Status": "RUNNING"
        },
        {
            "ID": "ASSERTION-InnerInnerDeployment-otherInnerInner()",
            "LogMessage": {},
            "Message": "[RUNNING] Assertion Test: InnerInnerDeployment.otherInnerInner()",
            "Name": "Assertion Test: InnerInnerDeployment.otherInnerInner()",
            "Status": "RUNNING"
        },
        {
            "ID": "ASSERTION-InnerDeployment-deployInnerInner()",
            "LogMessage": {},
            "Message": "[RUNNING] Assertion Test: InnerDeployment.deployInnerInner()",
            "Name": "Assertion Test: InnerDeployment.deployInnerInner()",
            "Status": "RUNNING"
        },
        {
            "ID": "ASSERTION-InnerDeployment-otherInner()",
            "LogMessage": {},
            "Message": "[RUNNING] Assertion Test: InnerDeployment.otherInner()",
            "Name": "Assertion Test: InnerDeployment.otherInner()",
            "Status": "RUNNING"
        },
        {
            "ID": "ASSERTION-InnerDeploymentFactory-deployInner()",
            "LogMessage": {},
            "Message": "[RUNNING] Assertion Test: InnerDeploymentFactory.deployInner()",
            "Name": "Assertion Test: InnerDeploymentFactory.deployInner()",
            "Status": "RUNNING"
        },
        {
            "ID": "ASSERTION-InnerDeploymentFactory-wee()",
            "LogMessage": {},
            "Message": "[RUNNING] Assertion Test: InnerDeploymentFactory.wee()",
            "Name": "Assertion Test: InnerDeploymentFactory.wee()",
            "Status": "RUNNING"
        }
    ]
}
```

- Logs info
```json
{
"logs": "⇾ Setting up base chain\n⇾ Initializing and validating corpus call sequences\n⇾ corpus: health: 77%, sequences: 1372 (1061 valid, 311 invalid)\n⇾ Fuzzing with 10 workers\n⇾ fuzz: elapsed: 0s, calls: 0 (0/sec), seq/s: 0, coverage: 1061\n⇾ fuzz: elapsed: 3s, calls: 23480 (7825/sec), seq/s: 103, coverage: 1061\n⇾ fuzz: elapsed: 6s, calls: 41060 (5858/sec), seq/s: 63, coverage: 1061\n⇾ fuzz: elapsed: 9s, calls: 58579 (5831/sec), seq/s: 61, coverage: 1061\n⇾ fuzz: elapsed: 12s, calls: 69870 (3688/sec), seq/s: 37, coverage: 1061\n⇾ fuzz: elapsed: 15s, calls: 80661 (3583/sec), seq/s: 36, coverage: 1061\n⇾ fuzz: elapsed: 18s, calls: 101069 (6778/sec), seq/s: 69, coverage: 1061\n⇾ fuzz: elapsed: 21s, calls: 115915 (4932/sec), seq/s: 49, coverage: 1061\n⇾ fuzz: elapsed: 24s, calls: 131499 (5186/sec), seq/s: 51, coverage: 1061\n"
}
```

- Coverage info
```json
{
    "Files": {
        "/home/aaliyah/Desktop/BRAINY_CODELAB/testing/nested_deployments.sol": {
            "Path": "/home/aaliyah/Desktop/BRAINY_CODELAB/testing/nested_deployments.sol",
            "Lines": [
                {
                    "IsActive": false,
                    "Start": 0,
                    "End": 38,
                    "Contents": "Ly8gU1BEWC1MaWNlbnNlLUlkZW50aWZpZXI6IFVubGljZW5zZQ==",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 38,
                    "End": 62,
                    "Contents": "cHJhZ21hIHNvbGlkaXR5IF4wLjguMDs=",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 62,
                    "End": 63,
                    "Contents": "",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 63,
                    "End": 95,
                    "Contents": "Y29udHJhY3QgSW5uZXJJbm5lckRlcGxveW1lbnQgew==",
                    "IsCovered": true,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 95,
                    "End": 135,
                    "Contents": "ICAgIGZ1bmN0aW9uIG90aGVySW5uZXJJbm5lcigpIHB1YmxpYyB7",
                    "IsCovered": true,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 135,
                    "End": 151,
                    "Contents": "ICAgICAgICByZXR1cm47",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 151,
                    "End": 157,
                    "Contents": "ICAgIH0=",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 157,
                    "End": 158,
                    "Contents": "",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 158,
                    "End": 234,
                    "Contents": "ICAgIGZ1bmN0aW9uIHByb3BlcnR5X2lubmVyX2lubmVyX2RlcGxveW1lbnQoKSBwdWJsaWMgdmlldyByZXR1cm5zIChib29sKSB7",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 234,
                    "End": 274,
                    "Contents": "ICAgICAgICAvLyBBU1NFUlRJT046IEZhaWwgaW1tZWRpYXRlbHku",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 274,
                    "End": 295,
                    "Contents": "ICAgICAgICByZXR1cm4gdHJ1ZTs=",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 295,
                    "End": 301,
                    "Contents": "ICAgIH0=",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 301,
                    "End": 302,
                    "Contents": "",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 302,
                    "End": 357,
                    "Contents": "ICAgIC8vIGZ1bmN0aW9uIGZ1enpfaW5uZXJfaW5uZXJfZGVwbG95bWVudCgpIHB1YmxpYyB7",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 357,
                    "End": 400,
                    "Contents": "ICAgIC8vICAgICAvLyBBU1NFUlRJT046IEZhaWwgaW1tZWRpYXRlbHku",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 400,
                    "End": 426,
                    "Contents": "ICAgIC8vICAgICBhc3NlcnQoZmFsc2UpOw==",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 426,
                    "End": 435,
                    "Contents": "ICAgIC8vIH0=",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 435,
                    "End": 437,
                    "Contents": "fQ==",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 437,
                    "End": 438,
                    "Contents": "",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 438,
                    "End": 465,
                    "Contents": "Y29udHJhY3QgSW5uZXJEZXBsb3ltZW50IHs=",
                    "IsCovered": true,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 465,
                    "End": 500,
                    "Contents": "ICAgIGZ1bmN0aW9uIG90aGVySW5uZXIoKSBwdWJsaWMgew==",
                    "IsCovered": true,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 500,
                    "End": 516,
                    "Contents": "ICAgICAgICByZXR1cm47",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 516,
                    "End": 522,
                    "Contents": "ICAgIH0=",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 522,
                    "End": 523,
                    "Contents": "",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 523,
                    "End": 593,
                    "Contents": "ICAgIGZ1bmN0aW9uIHByb3BlcnR5X2lubmVyX2RlcGxveW1lbnQoKSBwdWJsaWMgdmlldyByZXR1cm5zIChib29sKSB7",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 593,
                    "End": 633,
                    "Contents": "ICAgICAgICAvLyBBU1NFUlRJT046IEZhaWwgaW1tZWRpYXRlbHku",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 633,
                    "End": 654,
                    "Contents": "ICAgICAgICByZXR1cm4gdHJ1ZTs=",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 654,
                    "End": 660,
                    "Contents": "ICAgIH0=",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 660,
                    "End": 661,
                    "Contents": "",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 661,
                    "End": 720,
                    "Contents": "ICAgIGZ1bmN0aW9uIGRlcGxveUlubmVySW5uZXIoKSBwdWJsaWMgcmV0dXJucyAoYWRkcmVzcykgew==",
                    "IsCovered": true,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 720,
                    "End": 772,
                    "Contents": "ICAgICAgICByZXR1cm4gYWRkcmVzcyhuZXcgSW5uZXJJbm5lckRlcGxveW1lbnQoKSk7",
                    "IsCovered": true,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 772,
                    "End": 778,
                    "Contents": "ICAgIH0=",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 778,
                    "End": 780,
                    "Contents": "fQ==",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 780,
                    "End": 781,
                    "Contents": "",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 781,
                    "End": 848,
                    "Contents": "Ly8gVGVzdENvbnRyYWN0IGRlcGxveXMgSW5uZXJEZXBsb3ltZW50IHRvIHRlc3QgaW5uZXIgZGVwbG95bWVudHMu",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 848,
                    "End": 882,
                    "Contents": "Y29udHJhY3QgSW5uZXJEZXBsb3ltZW50RmFjdG9yeSB7",
                    "IsCovered": true,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 882,
                    "End": 910,
                    "Contents": "ICAgIGZ1bmN0aW9uIHdlZSgpIHB1YmxpYyB7",
                    "IsCovered": true,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 910,
                    "End": 926,
                    "Contents": "ICAgICAgICByZXR1cm47",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 926,
                    "End": 932,
                    "Contents": "ICAgIH0=",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 932,
                    "End": 933,
                    "Contents": "",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 933,
                    "End": 987,
                    "Contents": "ICAgIGZ1bmN0aW9uIGRlcGxveUlubmVyKCkgcHVibGljIHJldHVybnMgKGFkZHJlc3MpIHs=",
                    "IsCovered": true,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": true,
                    "Start": 987,
                    "End": 1034,
                    "Contents": "ICAgICAgICByZXR1cm4gYWRkcmVzcyhuZXcgSW5uZXJEZXBsb3ltZW50KCkpOw==",
                    "IsCovered": true,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 1034,
                    "End": 1040,
                    "Contents": "ICAgIH0=",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 1040,
                    "End": 1042,
                    "Contents": "fQ==",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                },
                {
                    "IsActive": false,
                    "Start": 1042,
                    "End": 1043,
                    "Contents": "",
                    "IsCovered": false,
                    "IsCoveredReverted": false
                }
            ]
        }
    }
}
```

- Corpus info
```json
{
    "unexecutedCallSequences": [
        [
            {
                "call": {
                    "from": "0x0000000000000000000000000000000000010000",
                    "to": "0xa647ff3c36cfab592509e13860ab8c4f28781a66",
                    "nonce": 0,
                    "value": "0x0",
                    "gasLimit": 12500000,
                    "gasPrice": "0x1",
                    "gasFeeCap": "0x0",
                    "gasTipCap": "0x0",
                    "data": "0x946c3724",
                    "dataAbiValues": {
                        "methodSignature": "deployInner()",
                        "inputValues": []
                    },
                    "AccessList": null,
                    "SkipAccountChecks": false
                },
                "blockNumberDelay": 0,
                "blockTimestampDelay": 0
            },
            {
                "call": {
                    "from": "0x0000000000000000000000000000000000030000",
                    "to": "0xa647ff3c36cfab592509e13860ab8c4f28781a66",
                    "nonce": 1,
                    "value": "0x0",
                    "gasLimit": 12500000,
                    "gasPrice": "0x1",
                    "gasFeeCap": "0x0",
                    "gasTipCap": "0x0",
                    "data": "0x946c3724",
                    "dataAbiValues": {
                        "methodSignature": "deployInner()",
                        "inputValues": []
                    },
                    "AccessList": null,
                    "SkipAccountChecks": false
                },
                "blockNumberDelay": 56076,
                "blockTimestampDelay": 360624
            },
            {
                "call": {
                    "from": "0x0000000000000000000000000000000000030000",
                    "to": "0x54919a19522ce7c842e25735a9cfecef1c0a06da",
                    "nonce": 2,
                    "value": "0x0",
                    "gasLimit": 12500000,
                    "gasPrice": "0x1",
                    "gasFeeCap": "0x0",
                    "gasTipCap": "0x0",
                    "data": "0xc23ac55e",
                    "dataAbiValues": {
                        "methodSignature": "otherInner()",
                        "inputValues": []
                    },
                    "AccessList": null,
                    "SkipAccountChecks": false
                },
                "blockNumberDelay": 1,
                "blockTimestampDelay": 134226
            },
            {
                "call": {
                    "from": "0x0000000000000000000000000000000000020000",
                    "to": "0xa647ff3c36cfab592509e13860ab8c4f28781a66",
                    "nonce": 0,
                    "value": "0x0",
                    "gasLimit": 12500000,
                    "gasPrice": "0x1",
                    "gasFeeCap": "0x0",
                    "gasTipCap": "0x0",
                    "data": "0x946c3724",
                    "dataAbiValues": {
                        "methodSignature": "deployInner()",
                        "inputValues": []
                    },
                    "AccessList": null,
                    "SkipAccountChecks": false
                },
                "blockNumberDelay": 39558,
                "blockTimestampDelay": 130342
            }
      ]
    ]
}
```