import { cubicIn } from 'svelte/easing';

export function shrink(node, params) {
    const existingTransform = getComputedStyle(node).transform.replace('none', '');

    return {
        delay: params.delay || 0,
        duration: params.duration || 400,
        easing: params.easing || cubicIn,
        css: (t, u) => `
            transform: ${existingTransform} scaleY(${t});
            height: ${t * node.offsetHeight}px;
            opacity: ${t};
            transform-origin: top;
        `
    };
}