<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\Api\ApiDocController;
use App\Http\Controllers\Api\ApiDocTreeController;
use App\Http\Controllers\Api\ApiUrlController;
use App\Http\Controllers\Api\DirectoryController;
use App\Http\Controllers\Api\ApiDocNoAuthController;
use App\Http\Controllers\Api\DocTrashController;
use App\Http\Controllers\Api\EmailController;
use App\Http\Controllers\Api\ProjectController;
use App\Http\Controllers\Api\ProjectGroupController;
use App\Http\Controllers\Api\ProjectMemberController;
use App\Http\Controllers\Api\ProjectParameterController;
use App\Http\Controllers\Api\ProjectsController;
use App\Http\Controllers\Api\ProjectNoAuthController;
use App\Http\Controllers\Api\TeamController;
use App\Http\Controllers\Api\UserController;
use App\Http\Controllers\Api\LogoutController;
use App\Http\Controllers\Api\UploadController;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| is assigned the "api" middleware group. Enjoy building your API!
|
*/

// 邮箱登录
Route::post('/email_login', [EmailController::class, 'login']);
// 邮箱注册
Route::post('/email_register', [EmailController::class, 'register']);
// 退出登录
Route::post('/logout', [LogoutController::class, 'index']);

// 个人设置
Route::prefix('user')->group(function () {
    // 获取个人信息
    Route::get('profile', [UserController::class, 'profile']);
    // 更换头像
    Route::post('change_avatar', [UserController::class, 'changeAvatar']);
    // 修改个人信息
    Route::post('change_profile', [UserController::class, 'changeProfile']);
    // 修改密码
    Route::post('change_password', [UserController::class, 'changePassword']);
});

// 团队成员
Route::prefix('team')->group(function () {
    // 成员列表
    Route::get('members', [TeamController::class, 'members']);
    // 成员详情
    Route::get('member', [TeamController::class, 'member']);
    // 成员参与的项目列表
    Route::get('member_projects', [TeamController::class, 'memberProjects']);
    // 添加成员
    Route::post('add_member', [TeamController::class, 'addMember']);
    // 修改成员信息
    Route::post('edit_member_info', [TeamController::class, 'editMemberInfo']);
    // 修改成员密码
    Route::post('edit_member_password', [TeamController::class, 'editMemberPassword']);
    // 删除成员
    Route::post('remove_member', [TeamController::class, 'removeMember']);
});

// 项目分组列表
Route::get('project_groups', [ProjectGroupController::class, 'projectGroups']);
// 项目分组
Route::prefix('project_group')->group(function () {
    // 添加分组
    Route::post('create', [ProjectGroupController::class, 'createGroup']);
    // 重命名分组
    Route::post('rename', [ProjectGroupController::class, 'changeName']);
    // 分组排序
    Route::post('change_order', [ProjectGroupController::class, 'changeOrder']);
    // 删除分组
    Route::post('remove', [ProjectGroupController::class, 'remove']);
});

// 项目列表
Route::prefix('projects')->group(function () {
    // 项目分组下的详细列表
    Route::get('/', [ProjectsController::class, 'index']);
    // 权限下的项目列表，只返回项目基本信息
    Route::get('/base', [ProjectsController::class, 'base']);
});

// 项目
Route::prefix('project')->group(function () {
    // 项目详情
    Route::get('/', [ProjectController::class, 'index']);
    // 项目状态
    Route::get('/status', [ProjectNoAuthController::class, 'status']);
    // 创建项目
    Route::post('create', [ProjectsController::class, 'create']);
    // 上传项目图标
    Route::post('icon', [ProjectsController::class, 'icon']);
    // 开启关闭项目分享
    Route::post('share', [ProjectController::class, 'share']);
    // 重置分享项目访问秘钥
    Route::post('reset_share_secretkey', [ProjectController::class, 'resetShareSecretKey']);
    // 项目设置
    Route::post('setting', [ProjectController::class, 'storeSetting']);
    // 项目分组
    Route::post('change_group', [ProjectController::class, 'changeGroup']);
    // 修改项目拥有者
    Route::post('change_owner', [ProjectController::class, 'changeOwner']);
    // 移交项目
    Route::post('transfer', [ProjectController::class, 'transfer']);
    // 退出项目
    Route::post('quit', [ProjectController::class, 'quit']);
    // 删除项目
    Route::post('remove', [ProjectController::class, 'remove']);

    // 项目成员分页列表
    Route::get('members', [ProjectMemberController::class, 'index']);
    // 项目成员全部列表(id name avatar)
    Route::get('member_userinfo_list', [ProjectMemberController::class, 'memberUserinfoList']);
    // 不在此项目的成员列表
    Route::get('without_members', [ProjectMemberController::class, 'notInProject']);
    // 添加成员
    Route::post('add_member', [ProjectMemberController::class, 'addMember']);
    // 修改成员权限
    Route::post('change_member_authority', [ProjectMemberController::class, 'changeMemberAuthority']);
    // 移出成员
    Route::post('remove_member', [ProjectMemberController::class, 'removeMember']);

    // 公共参数列表
    Route::get('params', [ProjectParameterController::class, 'index']);
    // 添加公共参数
    Route::post('add_param', [ProjectParameterController::class, 'create']);
    // 修改公共参数
    Route::post('edit_param', [ProjectParameterController::class, 'update']);
    // 删除公共参数
    Route::post('remove_param', [ProjectParameterController::class, 'remove']);
});

// 项目和单篇文档预览
Route::prefix('preview')->group(function () {
    // 项目信息
    // Route::get('project', [ProjectPreviewController::class, 'projectInfo']);
    // 获取api文档树
    // Route::get('api_nodes', [ProjectPreviewController::class, 'apiNodes']);
    // api文档详情
    // Route::get('api_doc', [ProjectPreviewController::class, 'apiDoc']);
    // 文档搜索
    // Route::get('search', [ProjectPreviewController::class, 'search']);
    // 私有项目秘钥校验
    // Route::post('check', [ProjectPreviewController::class, 'checkSecretKey']);
    // 文档详情
    // Route::get('single_doc', [DocPreviewController::class, 'doc']);
});

// 回收站
Route::prefix('doc')->group(function () {
    // 文档列表
    Route::get('trash', [DocTrashController::class, 'index']);
    // 恢复已删除的api文档
    Route::post('restore_api_doc', [DocTrashController::class, 'restoreApiDoc']);
});

// API文档树
Route::prefix('api_tree')->group(function () {
    // 所有分类和文档列表
    Route::get('/', [ApiDocTreeController::class, 'index']);
    // 节点排序
    Route::post('sort', [ApiDocTreeController::class, 'sort']);
});

// API常用url
Route::prefix('api_url')->group(function () {
    // 获取常用url
    Route::get('list', [ApiUrlController::class, 'list']);
    // 删除常用url
    Route::post('remove', [ApiUrlController::class, 'remove']);
});

// 目录
Route::prefix('dir')->group(function () {
    // 所有分类列表
    Route::get('list', [DirectoryController::class, 'list']);
    // 创建分类
    Route::post('create', [DirectoryController::class, 'store']);
    // 分类重命名
    Route::post('rename', [DirectoryController::class, 'rename']);
    // 删除分类
    Route::post('remove', [DirectoryController::class, 'remove']);
});

// API文档
Route::prefix('api_doc')->group(function () {
    // 创建文档
    Route::post('create', [ApiDocController::class, 'create']);
    // 创建HTTP API文档
    Route::post('http_template', [ApiDocController::class, 'httpTemplate']);
    // 编辑文档
    Route::post('update', [ApiDocController::class, 'update']);
    // 文档详情
    Route::get('/', [ApiDocController::class, 'detail']);
    // 文档重命名
    Route::post('rename', [ApiDocController::class, 'rename']);
    // 复制文档
    Route::post('copy', [ApiDocController::class, 'copy']);
    // 删除文档
    Route::post('remove', [ApiDocController::class, 'remove']);
    // 文档搜索
    Route::get('search', [ApiDocController::class, 'search']);
    // 文档分享详情
    Route::post('share_detail', [ApiDocController::class, 'shareDetail']);
    // 分享文档
    Route::post('share', [ApiDocController::class, 'share']);
    // 重置分享文档的访问密码
    Route::post('share_secretkey', [ApiDocController::class, 'shareSecretKey']);
    // 文档分享状态
    Route::get('has_shared', [ApiDocNoAuthController::class, 'hasShared']);
    // 文档秘钥校验
    Route::post('secretkey_check', [ApiDocNoAuthController::class, 'checkSecretKey']);
    // 文档导入
    Route::post('import', [ApiDocController::class, 'import']);
    // 文档导入结果查询
    Route::get('import_result', [ApiDocController::class, 'importResult']);
    // 文档导出
    Route::post('export', [ApiDocController::class, 'export']);
    // 文档导出结果查询
    Route::get('export_result', [ApiDocController::class, 'exportResult']);
});

// 块文件上传初始化
Route::post('/upload_init', [UploadController::class, 'init']);
// 上传块文件
Route::post('/upload_save', [UploadController::class, 'save']);
// 获取完整文件路径
Route::post('/upload_path', [UploadController::class, 'path']);
