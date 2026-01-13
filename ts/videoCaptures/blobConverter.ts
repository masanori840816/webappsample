export function toBase64(value: Uint8Array): string {
    const binaryString = Array.from(value)
        .map(byte => String.fromCharCode(byte))
        .join('');

    // Base64 encoding
    return btoa(binaryString);
}