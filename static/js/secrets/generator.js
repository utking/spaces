;(() => {
const genPass = (len, upper, nums, special) => {
    const lower = "abcdefghijklmnopqrstuvwxyz";
    const upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    const numChars = "0123456789";
    const specialChars = "!@#$%^&*()-_=+[]{}|;:,.<>?";
    let chars = lower;

    if (upper) chars += upperChars;
    if (nums) chars += numChars;
    if (special) chars += specialChars;

    let pass = "";
    for (let i = 0; i < len; i++) {
        const randIdx = Math.floor(Math.random() * chars.length);
        pass += chars[randIdx];
    }

    return pass;
}

const generate = () => {
    const len = parseInt(document.getElementById("len").value);
    const upper = document.getElementById("upper").checked;
    const nums = document.getElementById("nums").checked;
    const special = document.getElementById("special").checked;

    const pass = genPass(len, upper, nums, special);
    document.getElementById("passOut").value = pass;
}

const reset = () => {
    document.getElementById("len").value = 16;
    document.getElementById("upper").checked = true;
    document.getElementById("nums").checked = true;
    document.getElementById("special").checked = true;
    document.getElementById("passOut").value = "";
}

document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("generate").addEventListener("click", generate);
    document.getElementById("reset").addEventListener("click", reset);
});
})();