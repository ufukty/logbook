/**
* @param {HTMLElement} element
*/
function removeAllChildren(element) {
    while (element.firstChild) {
        element.removeChild(element.firstChild);
    }
}

/***
* @param {Object} qr
* @param {[string]} qr.QrMatrix
* @param {string} qr.Text
*/
function updateUI(qr) {
    let qrContainer = document.getElementById("qr-container")
    removeAllChildren(qrContainer)
    qrContainer.style.setProperty("--length", `${qr.QrMatrix.length}`) // square
    for (const line of qr.QrMatrix) {
        for (const char of line) {
            var d = document.createElement("div")
            if (char === 'X') {
                d.classList.add("data")
            }
            qrContainer.appendChild(d)
        }
    }

    let secretContainer = document.getElementById("totp-secret")
    secretContainer.innerText = qr.Text
}

document.addEventListener("DOMContentLoaded", () => {
    let qr = {
        "QrMatrix": [
            "XXXXXXX XXX XXXXX  X X XX XX X  X   XX   XX XXXX  XXXXXXX",
            "X     X XXXXX X  X   XXXX XX         XX  XX XX X  X     X",
            "X XXX X X X X XX X   XXX  XX X  X X X   X   X XX  X XXX X",
            "X XXX X  XX XXX X XX  X X XX X  XX  X    XX  X X  X XXX X",
            "X XXX X    X  XXX  XX  XXXXXXXX    X X XXXX X  X  X XXX X",
            "X     X X     X X  X XXX  X   X XX XXX   X  X X   X     X",
            "XXXXXXX X X X X X X X X X X X X X X X X X X X X X XXXXXXX",
            "        XX X   XXXX XX   XX   X X  X XXXX X     X        ",
            "  XXX X X    X X XX   X  XXXXXXX   XX X  X XXXXX XXX  XXX",
            "    X      X X    X X X   X X  XX   XXXXX X X   XX  X X  ",
            "X X XXXXXXX X   X X X X  XXX XX  XX X    X XXX X  X XX X ",
            "X XX X X   X X XX  X  X  XX   XX  X  XX XX  X   XXX XXX  ",
            "  XXXXXX XXXX  XXXXX XXXXX X  X XX X XXX X X     XX  X  X",
            "XX XXX   X XXX XX  X XX X      XX    X   XX    XXX X  X X",
            "   X  X    X X X  XX  X   XX  X X   XXX X     X X X    X ",
            "XXX XX   XX X  X   XXX XX XX XXXX    X  XXX X XX X   XXXX",
            " X    XXXXX  X  XXXXX   X   XXX  X  X    XXX  X   XXX    ",
            "   X X X X  X XX  XXXXXX XX   X X X XX  XXXXX  XXX XX X X",
            "XXX X X  XX XX X X    XXX   X X XXXXXX X     XXXX X   X  ",
            " X XX   X XXX  XXXX  X X  X XX  X XX XXXX X X        XX  ",
            "XXX XXXX  XXX XXXX X XX X      XXX XXX     X  X   XX  X  ",
            " XX X   XX X  X X  XXXX  XX XX X    XX X    X  XX      X ",
            "XX  XXXXX XX  X  XXX XX X X XXX XX X X   X   XX XXXXX  X ",
            "XXX X   XXX     XX  X  XXXXXX XXX XX  X   XXX  X X XXXX  ",
            " XX  XX XXX X X  XX  X X X XXXXXX  X X X XXX  X       X  ",
            "  X  X X  X  X  XXX  XX XXXXXXXXX    X    X    XX   X    ",
            "X XXXXXXXXX XXX  XX  X  X XXXXX XX   XXX XXX XX XXXXX X  ",
            "X X X   X    X  X XX    X X   XXXXXXX  XXX  X XXX   XXXX ",
            " XXXX X X XXXXX  X X  X X X X X          X X  X X X X  X ",
            "XX  X   XXX XXXX XX X    XX   XX     X  X X X   X   XX  X",
            "   XXXXXXX  XXXXXXX  X   XXXXXXX   XX XXX    XXXXXXXXX   ",
            " X  X  XXXX X XXXXX   X  XXXXX  XXXXXX XXXXXXX X XX XXXX ",
            "  X  XX  XXX XXXXX  XXXXXXXXXX  X XXX     X  X    X X XXX",
            "X    X XX  X  XXXXXX XX  XXX     XX X  XXXXXXX  X    XXXX",
            "    X X X  XX  XXX XXXX  X X  XXXXX XXXX    XXX  X   XXX ",
            "XX   X XX X  X   XXXX    X  XX  XX     XXXXXXXX      XXX ",
            " X X  XXX XXX    XXX X    XXXX    X XX X  X   X X XXXX X ",
            " X   X  XX  XX X  X  XXX XX XXXXX XXXXX XXX   XX X X  XXX",
            "X  XXXXX  X XX     XX X    X  XX  XXXX XX   XX X XX  XXX ",
            " X   X X XXX  X XX  XXXX XXXX  XX X    XXX       XXX  XXX",
            "XX    X  X XXX XX XXX XXXX X   X X   X X XX  XXX    XX   ",
            "   X X X XXX X   X     X  X   X  XXX   XX X   XXX  X XX  ",
            "X  X XXXX  X XXXX XX XXXX XXX   X XX X   X XXXX  X XXXXX ",
            " X   X X XX  X XXXXXX   XX      XX X  XXXX X   XX XX XXXX",
            " X    X X  X  X XX       X X X  XX  X    X XXX   X XX    ",
            "  X XX    X   X X  XX X  X XX X  XXXX XX XXX     XX    XX",
            "X X  XX X X X  X   XX  XX  X  XX XX X XXX  XXXXXXX X XXX ",
            "XXXXX  X  X X  X XX  X      XX     X XX X XXX X   X XXXX ",
            "      X X X X XXXXX XXXX XXXXXX   XXX X  X X  X XXXXX  X ",
            "        X XXXXX  XX  X XX X   XX XX  XX  XX X  XX   X XXX",
            "XXXXXXX     XX XX       XXX X XXXX  X    X X XX X X XXXX ",
            "X     X   X    X   X  X   X   X    X    XXX X XXX   X X  ",
            "X XXX X X   XXXX  X XX X XXXXXXXXXXXX      X  XXXXXXX  X ",
            "X XXX X XXX XXXXX X XX XX  X  X   XX XX X XXX       X X  ",
            "X XXX X XX X    XXX  XXXX XX X  X  XX   X  X  XXXX   XX  ",
            "X     X   X XXXX    X XXXXX XX X XXX X X X XXXX XX  X X  ",
            "XXXXXXX  X X XX X XXX XXXX  X X X    XX X  X   XX X XXXX "
        ],
        "Text": "QIHBLVXHJT4HJBWG7CHNSUBPW6MUUUJV"
    }
    updateUI(qr)
})