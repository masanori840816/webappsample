/**
 * Remove video codecs from SDP except specified ones
 * resource: https://stackoverflow.com/questions/52738290/how-to-remove-video-codecs-in-webrtc-sdp
 * @param orgsdp 
 * @param codec 
 * @returns 
 */
export function removeVideoCodec(orgsdp: string, codec: string) {
    const internalFunc = (sdp: string): string => {
        const codecre = new RegExp('(a=rtpmap:(\\d*) ' + codec + '/90000\\r\\n)');
        const rtpmaps = sdp.match(codecre);        
        if (rtpmaps == null || rtpmaps.length <= 2) {
            return sdp;
        }
        const rtpmap = rtpmaps[2];
        let modsdp = sdp.replace(codecre, "");
        modsdp = modsdp.replace(new RegExp('(a=rtcp-fb:' + rtpmap + '.*\r\n)', 'g'), "");
        modsdp = modsdp.replace(new RegExp('(a=fmtp:' + rtpmap + '.*\r\n)', 'g'), "");
        const aptpre = new RegExp('(a=fmtp:(\\d*) apt=' + rtpmap + '\\r\\n)');
        const aptmaps = modsdp.match(aptpre);
        let fmtpmap = "";
        if (aptmaps != null && aptmaps.length >= 3) {
            fmtpmap = aptmaps[2] ?? "";
            modsdp = modsdp.replace(aptpre, "");
            // remove video codecs by the "fmtpmap" value.
            modsdp = modsdp.replace(new RegExp('(a=rtpmap:' + fmtpmap + '.*\r\n)', 'g'), "");
            modsdp = modsdp.replace(new RegExp('(a=rtcp-fb:' + fmtpmap + '.*\r\n)', 'g'), "");
        }
        const videore = /(m=video.*\r\n)/;
        const videolines = modsdp.match(videore);
        if (videolines != null) {
            //If many m=video are found in SDP, this program doesn't work.
            const videoline = videolines[0].substring(0, videolines[0].length - 2);
            const videoelem = videoline.split(" ");
            let modvideoline = videoelem[0] ?? "";
            for (let i = 1; i < videoelem.length; i++) {
                if (videoelem[i] == rtpmap || videoelem[i] == fmtpmap) {
                    continue;
                }
                modvideoline += " " + videoelem[i];
            }
            modvideoline += "\r\n";
            modsdp = modsdp.replace(videore, modvideoline);
        }
        return internalFunc(modsdp);
    };
    return internalFunc(orgsdp);
}