const setTimeoutPromise = ms => new Promise(resolve => setTimeout(resolve, ms));

export const copyToClipBoard = async text => {
    const el = document.createElement("textarea");
    el.value = text;
    el.setAttribute("readonly", "");
    el.style.position = "absolute";
    el.style.left = "-9999px";
    document.body.appendChild(el);
    await setTimeoutPromise(50);
    el.select();
    document.execCommand("copy");
    document.body.removeChild(el);
};