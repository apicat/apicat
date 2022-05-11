<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\Mock\MockController;
use App\Http\Controllers\Mock\ImageController;
use App\Http\Controllers\Mock\FileController;

/**
 * Mock API Routes
 */
Route::get('/mock_image/{imgInfo}', [ImageController::class, 'index']);
Route::get('/mock_file/{fileInfo}', [FileController::class, 'index']);
Route::any('/{projectId}/{path}', [MockController::class, 'index'])->where([
    'projectId' => '[0-9]+',
    'path' => '.*'
]);
