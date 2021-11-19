if (typeof (Number.prototype.toRadians) === "undefined") {
    Number.prototype.toRadians = function () {
        return this * Math.PI / 180;
    }        
}

/**
 * Calculate distance between two geographical coordinates in earth surface.
 * 
 * @param {number} x1 latitude for the first coordinate.
 * @param {number} y1 longitute for the first coordinate.
 * @param {number} x2 latitude for the second coordinate.
 * @param {number} y2 longitude for the second coordinate.
 * 
 * @return {number} Distance in kilometers.
 * 
 */

const haversine(x1, y1, x2, y2) = () => {    
    const R = 6371;
    const dLat = (x2 - x1).toRadians();
    const dLng = (y2 - y1).toRadians();
    const a =
        Math.pow(Math.sin(dLat / 2), 2) +
        Math.pow(Math.sin(dLng / 2), 2) *
        Math.cos((x1).toRadians()) * Math.cos((x2).toRadians());
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    const d = R * c;
    return d;
}
