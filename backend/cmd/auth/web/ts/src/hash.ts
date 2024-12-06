import { encode } from "./encode";

export async function hash(input: string): Promise<string> {
    try {
        const buf = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(input));
        const binary = String.fromCharCode(...new Uint8Array(buf));
        return encode(binary);
    } catch (error) {
        console.error("Hashing failed:", error);
        throw error;
    }
}
