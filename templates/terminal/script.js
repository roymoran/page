document.addEventListener('DOMContentLoaded', function () {
    const content = document.querySelector('.content');
    const initialPrompt = getInitialPrompt();


    // List of commands to simulate with their corresponding multiple outputs
    const commands = [
        {
            command: 'page new', outputs: [
                { text: '', delay: 0, animationSequence: { loading: [], complete: "" } },
            ]
        },
        {
            command: 'page up', outputs: [
                { text: 'Checking Host............', delay: 8000, animationSequence: { loading: ['Provisioning...[-]', 'Provisioning...[\\]', 'Provisioning...[|]', 'Provisioning...[/]'], complete: "Provisioned...[✓]" } },
                { text: 'Checking Certificate.....', delay: 5000, animationSequence: { loading: ['Generating....[-]', 'Generating....[\\]', 'Generating....[|]', 'Generating....[/]'], complete: "Generated.....[✓]" } },
                { text: 'Checking Website Files...', delay: 6000, animationSequence: { loading: ['Uploading.....[-]', 'Uploading.....[\\]', 'Uploading.....[|]', 'Uploading.....[/]'], complete: "Uploaded......[✓]" } },
                { text: 'Checking Domain..........', delay: 6000, animationSequence: { loading: ['Updating......[-]', 'Updating......[\\]', 'Updating......[|]', 'Updating......[/]'], complete: "Updated.......[✓]" } },
                { text: '\n', delay: 0, animationSequence: { loading: [], complete: "" } },
                { text: 'Page details', delay: 0, animationSequence: { loading: [], complete: "" } },
                { text: 'Domain: https://pagecli.com', delay: 0, animationSequence: { loading: [], complete: "" } },
                { text: `Certificate Expires: ${getFormattedDate()}`, delay: 0, animationSequence: { loading: [], complete: "" } },
                { text: 'Certificate Renewal: run page up in 60 days', delay: 0, animationSequence: { loading: [], complete: "" } },
            ]
        },
        {
            // empty command to simulate a blank line
            command: '', outputs: [
                { text: '', delay: 0, animationSequence: { loading: [], complete: "" } },
            ]
        }
    ];

    let commandIndex = 0;
    let typingSpeed = 100; // speed of typing in milliseconds
    let resetAnimationDelay = 8000; // time to wait before resetting the whole animation


    function getLastLoginText() {
        const currentDate = new Date();
        const dayOfWeek = currentDate.toLocaleString('en-US', { weekday: 'short' }); // "Wed"
        const month = currentDate.toLocaleString('en-US', { month: 'short' }); // "Nov"
        const dayOfMonth = currentDate.getDate(); // 15
        const year = currentDate.getFullYear(); // 2023
        const time = currentDate.toTimeString().split(' ')[0]; // "08:33:20"
        return `Last login: ${dayOfWeek} ${month} ${dayOfMonth} ${year} ${time} on ttys000`;
    }

    function resetTerminal(resetDelay) {
        // Clear the terminal after a delay and then restart the animation
        setTimeout(() => {
            const loginText = getLastLoginText();
            content.innerHTML = `<p>${loginText}</p>`;
            commandIndex = 0;
            typeCommand(commands[commandIndex]);
        }, resetDelay);
    }

    function displayOutput(commandObj, outputIndex = 0) {
        if (outputIndex < commandObj.outputs.length) {
            const outputObj = commandObj.outputs[outputIndex];
            if (outputObj.animationSequence.loading.length > 0) {
                displayAnimatedOutput(outputObj, () => displayOutput(commandObj, outputIndex + 1));
            } else {
                const output = document.createElement('p');
                output.textContent = outputObj.text;
                content.appendChild(output);
                scrollToBottom();

                setTimeout(() => {
                    displayOutput(commandObj, outputIndex + 1);
                }, outputObj.delay);
            }
        } else {
            // After all outputs for this command, move to the next command or reset
            if (commandIndex < commands.length - 1) {
                commandIndex++;
                typeCommand(commands[commandIndex]);
            } else {
                // All commands executed, reset
                resetTerminal(resetAnimationDelay);
            }
        }
    }

    function typeCommand(commandObj) {
        let index = 0;
        const line = document.createElement('p');
        const prompt = document.createElement('span');
        prompt.textContent = initialPrompt;
        const commandSpan = document.createElement('span');
        line.appendChild(prompt);
        line.appendChild(commandSpan);
        content.appendChild(line);

        function type() {
            if (index < commandObj.command.length) {
                commandSpan.textContent += commandObj.command[index++];
                setTimeout(type, typingSpeed);
            } else {
                // After finishing typing, display the output and prepare for next command
                displayOutput(commandObj);
                scrollToBottom();
            }
        }
        type();
    }

    function scrollToBottom() {
        content.scrollTop = content.scrollHeight;
    }

    // Function to determine the initial prompt based on the user's OS
    function getInitialPrompt() {
        const userAgent = navigator.userAgent.toLowerCase();
        if (userAgent.includes('mac')) {
            return 'developer@mac ~ % ';
        } else if (userAgent.includes('win')) {
            return 'C:\\Users\\developer> ';
        } else if (userAgent.includes('linux')) {
            return 'developer@linux ~ $ ';
        } else {
            return 'developer@system ~ % ';
        }
    }

    function displayAnimatedOutput(outputObj, callback) {
        const output = document.createElement('p');
        output.innerHTML = outputObj.text + outputObj.animationSequence.loading[0];
        content.appendChild(output);
        scrollToBottom();

        // Use outputObj.delay as the total duration for the animation
        animate(output, outputObj.text, outputObj.animationSequence, outputObj.delay, callback);
    }

    function animate(element, baseText, animations, totalDuration, callback) {
        let index = 0;
        const animationInterval = 200; // constant interval for animation frames in milliseconds
        const totalAnimationCycles = Math.floor(totalDuration / animationInterval);
        let currentCycle = 0;

        function updateAnimation() {
            if (index >= animations.loading.length) {
                index = 0; // Reset the index to loop the animation
            }
            element.innerHTML = baseText + animations.loading[index++];
            currentCycle++;

            // Continue the animation or end it based on the total duration
            if (currentCycle < totalAnimationCycles) {
                setTimeout(updateAnimation, animationInterval);
            } else {
                element.innerHTML = baseText + animations.complete;
                callback();
            }
        }

        updateAnimation();
    }

    function getFormattedDate() {
        const date = new Date();
        date.setDate(date.getDate() + 90); // Add 90 days to the current date
        const options = { 
            year: 'numeric', 
            month: 'long', 
            day: 'numeric', 
            hour: '2-digit', 
            minute: '2-digit', 
            hour12: true,
            timeZoneName: 'short' 
        };
        return date.toLocaleString('en-US', options).replace(/,/g, '');
    }
    
    // Call resetTerminal initially to set the date and start the first command
    resetTerminal(0);
});
