<?php

namespace App\Http\Controllers\Api;

use App\Exceptions\NotFoundException;
use App\Exceptions\SecretKeyExpiredException;
use App\Http\Controllers\Controller;
use App\Models\Project;
use App\Repositories\Project\ApiDocRepository;
use App\Repositories\Project\ProjectMemberRepository;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\Project\ProjectShareRepository;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Str;
use Illuminate\Validation\ValidationException;
use App\Modules\EditorJsonToHtml\Parser;
use Psr\SimpleCache\InvalidArgumentException;

class ProjectPreviewController extends Controller
{
    /**
     * 项目信息
     * @param Request $request
     * @return array
     * @throws NotFoundException
     */
    public function projectInfo(Request $request)
    {
        $request->validate([
            'project_id' => 'required|integer|min:1'
        ]);

        $project = ProjectRepository::get($request->input('project_id'));
        if (!$project) {
            throw new NotFoundException;
        }

        if (Auth::guard('api')->check()) {
            // 用户处于登录状态
            $inThisProject = ProjectMemberRepository::inThisProject($project->id, Auth::guard('api')->id());
        } else {
            $inThisProject = false;
        }

        // 如果项目是公开项目，默认为true，否则去检查项目是否被分享过
        $hasShared = $project->visibility || ProjectShareRepository::hasShared($project->id);

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'id' => $project->id,
                'name' => $project->name,
                'visibility' => $project->visibility,
                'in_this' => $inThisProject,
                'has_shared' => $hasShared
            ]
        ];
    }

    /**
     * 获取api文档树
     * @param Request $request
     * @return array
     * @throws SecretKeyExpiredException
     * @throws ValidationException
     * @throws InvalidArgumentException
     */
    public function apiNodes(Request $request)
    {
        $request->validate([
            'project_id' => 'required|integer|min:1',
            'token' => 'nullable|string|size:60'
        ]);

        $project = $this->updateCacheLifeTime($request);

        $tree = ApiDocRepository::getTree($project->id);
        return ['status' => 0, 'msg' => '', 'data' => $tree];
    }
}
