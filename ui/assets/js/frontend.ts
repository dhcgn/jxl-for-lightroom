// @ts-check

let stopGetProgress = false
async function getProgress() {
    if (stopGetProgress)
        return

    try {
        // 👇️ const response: Response
        const response = await fetch('/progress', {
            method: 'GET',
            headers: {
                Accept: 'application/text',
            },
        });

        if (!response.ok) {
            throw new Error(`Error! status: ${response.status}`);
        }

        // 👇️ const result: GetUsersResponse
        const result = await response.text();

        console.log('result is: ', result);

        return result;
    } catch (error) {
        stopGetProgress = true
        if (error instanceof Error) {
            console.log('error message: ', error.message);
            return error.message;
        } else {
            console.log('unexpected error: ', error);
            return 'An unexpected error occurred';
        }
    }
}

const sleep = ms => {
    return new Promise(resolve => setTimeout(resolve, ms))
}

const app = document.getElementById("progressbarts");

let lastProgress = "-1"

async function updateProgress() {
    const progress = sleep(100).then(n=>getProgress())
    let p = await progress
    if (lastProgress == p)
        return

    lastProgress = p

    if (app) {
       
        app.ariaValueNow = p
        app.innerText = p
        app.style.width = `${p}%`
        console.log('set dom ');
    }
}


// for (; ;) {
//     updateProgress()
// }

(async () => {
    try {
        for (; ;) {
            
            // await sleep(100);
            await updateProgress()
        }
    } catch (e) {
        // Deal with the fact the chain failed
    }
    // `text` is not available here
})();