<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;
use App\Repositories\DownloadRepository;

class DownloadController extends Controller
{
    public function index($fileName)
    {
        $downloadRepository = new DownloadRepository();
        if (!$downloadRepository->check($fileName)) {
            abort(404);
        }

        return Storage::disk('export')->download($fileName, $downloadRepository->downloadName());
    }
}
