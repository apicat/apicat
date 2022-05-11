<?php

namespace App\Http\Controllers\Mock;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use App\Exceptions\NotFoundException;
use App\Modules\Mock\Mocks\Base\FileRandom;

class FileController extends Controller
{
    protected $types = [
        'md' => 'md',
        'docx' => 'word',
        'xlsx' => 'excel',
        'csv' => 'csv'
    ];

    public function index(Request $request, $fileInfo)
    {
        if (strpos($fileInfo, '.') === false) {
            throw new NotFoundException;
        }

        list($fileName, $type) = explode('.', $fileInfo, 2);

        if (!isset($this->types[$type])) {
            throw new NotFoundException;
        }

        if (!$file = FileRandom::path($this->types[$type])) {
            throw new NotFoundException;
        }

        if ($file['name'] != $fileInfo) {
            throw new NotFoundException;
        }

        return response()->file($file['path']);
    }
}
