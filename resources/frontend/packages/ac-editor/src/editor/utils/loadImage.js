const IMAGE_CACHE = {};
window.IMAGE_CACHE = IMAGE_CACHE;

export const loadImage = (src) =>
  new Promise((resolve, reject) => {
    const result = {
      complete: false,
      width: 0,
      height: 0,
      src,
    };

    if (!src) {
      reject(result);
      return;
    }

    if (IMAGE_CACHE[src]) {
      resolve({ ...IMAGE_CACHE[src] });
      return;
    }

    const img = new Image();

    img.onload = () => {
      result.width = img.width;
      result.height = img.height;
      result.complete = true;

      IMAGE_CACHE[src] = { ...result };
      resolve(result);
    };

    img.onerror = () => {
      reject(result);
    };

    img.src = src;
  });

export default loadImage;
