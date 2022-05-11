<?php

namespace App\Http\Controllers\Mock;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use App\Exceptions\NotFoundException;
use App\Modules\Mock\Mocks\Image;

class ImageController extends Controller
{
    protected $types = [
        'jpg' => 'image/jpeg',
        'jpeg' => 'image/jpeg',
        'png' => 'image/png',
    ];

    public function index(Request $request, $imgInfo)
    {
        if (strpos($imgInfo, '.') === false) {
            throw new NotFoundException;
        }

        list($size, $type) = explode('.', $imgInfo, 2);

        if (!isset($this->types[$type])) {
            throw new NotFoundException;
        }

        if (strpos($imgInfo, 'x') === false) {
            throw new NotFoundException;
        }

        list($width, $height) = explode('x', $size, 2);
        $width = (int)$width;
        $height = (int)$height;

        if ($width < 20 or $width > 1024) {
            throw new NotFoundException;
        }
        if ($height < 20 or $height > 1024) {
            throw new NotFoundException;
        }

        $response = response()->make(Image::stream($width, $height, $type));
        $response->header('Content-Type', $this->types[$type]);
        return $response;
    }
}
