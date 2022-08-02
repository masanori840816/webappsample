export function getBoolParam(key: string): boolean {
    const params = new URLSearchParams(location.search);
    const target = params.get(key);
    if (target == null || target.length <= 0) {
        return false;
    }
    return target === "true";
}