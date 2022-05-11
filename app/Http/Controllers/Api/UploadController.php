<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Validation\ValidationException;
use App\Repositories\FileChunkUploadRepository;

class UploadController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api']);
    }

    public function init(Request $request)
    {
        $request->validate([
            'name' => 'required|string|min:4',
            'chunks' => 'required|integer|min:1',
            'fileSize' => 'required|integer|max:10240'
        ]);

        $fileChunkUploadRepository = new FileChunkUploadRepository;

        if (!$fileChunkUploadRepository->checkAllow($request->input('name'))) {
            throw ValidationException::withMessages([
                'name' => '不支持该文件类型',
            ]);
        }

        if (!$fileChunkUploadRepository->checkChunks($request->input('fileSize'), $request->input('chunks'))) {
            throw ValidationException::withMessages([
                'chunks' => '文件块数量不正确',
            ]);
        }

        if (!$fileChunkUploadRepository->checkFileSize($request->input('fileSize'))) {
            throw ValidationException::withMessages([
                'chunks' => '文件过大',
            ]);
        }

        if (!$fileID = $fileChunkUploadRepository->generateCache(Auth::id(), $request->input('name'))) {
            throw ValidationException::withMessages([
                'result' => '上传失败，请稍后重试。',
            ]);
        }
        return response()->json(['status' => 0, 'msg' => '', 'data' => ['file_id' => $fileID]]);
    }

    public function save(Request $request)
    {
        $request->validate([
            'file_id' => 'required|string|size:32',
            'chunk_id' => 'required|integer|min:0',
            'file' => 'required|file|max:2048'
        ]);

        $fileChunkUploadRepository = new FileChunkUploadRepository;
        $userID = $fileChunkUploadRepository->getCache($request->input('file_id'));
        if (!$userID or $userID != Auth::id()) {
            throw ValidationException::withMessages([
                'result' => '上传失败，请稍后重试。',
            ]);
        }

        if (!$fileChunkUploadRepository->saveChunk($request->input('chunk_id'), $request->file('file')->get())) {
            throw ValidationException::withMessages([
                'result' => '上传失败，请稍后重试。',
            ]);
        }
        return response()->json(['status' => 0, 'msg' => '']);
    }

    public function path(Request $request)
    {
        $request->validate([
            'file_id' => 'required|string|size:32'
        ]);

        $fileChunkUploadRepository = new FileChunkUploadRepository;
        $userID = $fileChunkUploadRepository->getCache($request->input('file_id'));
        if (!$userID or $userID != Auth::id()) {
            throw ValidationException::withMessages([
                'result' => '文件不存在',
            ]);
        }

        if (!$path = $fileChunkUploadRepository->mergeChunk()) {
            throw ValidationException::withMessages([
                'result' => '上传失败，请稍后重试。',
            ]);
        }
        return response()->json(['status' => 0, 'msg' => '', 'data' => $path]);
    }
}
