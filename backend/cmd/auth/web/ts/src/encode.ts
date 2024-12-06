export function encode(input: string): string {
    return window.btoa(input).replaceAll("=", "");
}
