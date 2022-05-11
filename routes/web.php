<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\DownloadController;

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

// 单个文档预览
Route::view('doc/{docID}', 'index')->name('doc.index');
Route::view('doc/{docID}/verification', 'index')->name('doc.verification');

// 下载文件
Route::get('/download/{fileName}', [DownloadController::class, 'index'])->name('download');

Route::fallback(function () {
    return view('index');
});
