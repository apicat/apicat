<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\DownloadController;
use App\Http\Controllers\HomeController;

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
*/

Route::view('/', 'index')->name('welcome');

// 编辑项目
Route::view('editor/{projectID}/doc', 'index')->name('editor.doc');
Route::view('editor/{projectID}/doc/{nodeID}', 'index')->name('editor.detail');
Route::view('editor/{projectID}/doc/{nodeID}/edit', 'index')->name('editor.doc_edit');

// 项目预览
Route::view('app/{projectID}', 'index')->name('app.index');
Route::view('app/{projectID}/verification', 'index')->name('app.verification');
Route::view('app/{projectID}/{docID}', 'index')->name('app.detail');

// 项目列表
Route::get('/projects', [HomeController::class, 'index'])->name('projects');

// 项目
Route::prefix('project')->name('project.')->group(function () {
    // 编辑项目
    Route::get('{projectID}', [HomeController::class, 'project'])->name('index');
    Route::get('{projectID}/doc', [HomeController::class, 'project']);
    Route::get('{projectID}/doc/{docID}', [HomeController::class, 'project'])->name('doc');
    Route::get('{projectID}/doc/{docID}/edit', [HomeController::class, 'project']);
    Route::get('{projectID}/{any}', [HomeController::class, 'index']);
});

// 迭代列表
Route::get('/iterations', [HomeController::class, 'index'])->name('iterations');

// 迭代
Route::prefix('iteration')->name('iteration.')->group(function () {
    // 编辑项目
    Route::get('{iterationID}', [HomeController::class, 'iteration'])->name('index');
    Route::get('{iterationID}/doc', [HomeController::class, 'iteration']);
    Route::get('{iterationID}/doc/{docID}', [HomeController::class, 'iteration'])->name('doc');
    Route::get('{iterationID}/doc/{docID}/edit', [HomeController::class, 'iteration']);
    Route::get('{iterationID}/{any}', [HomeController::class, 'index']);
});

// 已删除api文档的预览
Route::get('/doc/{projectID}/trash_api_preview/{docID}', [HomeController::class, 'doc']);

// 单个文档预览
Route::prefix('doc')->name('doc.')->group(function () {
    Route::get('{docID}/verification', [HomeController::class, 'shareDoc'])->name('doc.verification');
    Route::get('{docID}', [HomeController::class, 'shareDoc'])->name('index');
});

// 下载文件
Route::get('/download/{fileName}', [DownloadController::class, 'index'])->name('download');

Route::fallback(function () {
    return view('index');
});
