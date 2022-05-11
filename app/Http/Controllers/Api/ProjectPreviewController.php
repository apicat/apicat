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

    /**
     * api文档详情
     * @param Request $request
     * @return array
     * @throws NotFoundException
     * @throws SecretKeyExpiredException
     * @throws ValidationException
     */
    public function apiDoc(Request $request)
    {
        $request->validate([
            'project_id' => 'required|integer|min:1',
            'token' => 'nullable|string|size:60',
            'node_id' => 'required|integer|min:1'
        ]);

        $project = $this->updateCacheLifeTime($request);

        $node = ApiDocRepository::getNode($request->input('node_id'));
        if (!$node or $node->project_id != $project->id) {
            throw new NotFoundException;
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'id' => $node->id,
                'project_id' => $node->project_id,
                'title' => $node->title,
                'content' => $node->content ? Parser::parse($node->content, $project->id, $node->id) : '',
                'updated_time' => $node->updated_at->format('Y-m-d H:i')
            ]
        ];
    }

    /**
     * 文档搜索
     * @param Request $request
     * @return array
     * @throws SecretKeyExpiredException
     * @throws ValidationException
     */
    public function search(Request $request)
    {
        $request->validate([
            'project_id' => 'required|integer|min:1',
            'token' => 'nullable|string|size:60',
            'keywords' => 'required|string|max:255'
        ]);

        $project = $this->updateCacheLifeTime($request);

        $records = ApiDocRepository::searchNode($project->id, $request->input('keywords'));
        if ($records) {
            foreach ($records as $k => $v) {
                $records[$k]['link'] = route('app.detail', ['projectID' => $project->id, 'docID' => $v['doc_id']]);
            }
        }

        return ['status' => 0, 'msg' => '', 'data' => $records];
    }

    /**
     * 私有项目秘钥校验
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function checkSecretKey(Request $request)
    {
        $request->validate([
            'project_id' => 'required|integer|min:1',
            'secret_key' => 'required|string|size:6'
        ]);

        $project = ProjectRepository::get($request->input('project_id'));

        if (!ProjectShareRepository::check($project->id, $request->input('secret_key'))) {
            throw ValidationException::withMessages([
                'secret_key' => '访问密码不正确',
            ]);
        }

        $token = Str::random(60);
        $storageKey = hash('sha256', $token);

        if (!Cache::put($storageKey, ['project_id' => $project->id], 7200)) {
            throw ValidationException::withMessages([
                'result' => '验证失败，请稍后重试。',
            ]);
        }

        return ['status' => 0, 'msg' => '', 'data' => $token];
    }

    /**
     * 更新缓存生命周期
     * @param Request $request
     * @return Project
     * @throws SecretKeyExpiredException
     * @throws ValidationException
     */
    protected function updateCacheLifeTime($request)
    {
        $project = ProjectRepository::get($request->input('project_id'));

        if ($project->visibility == 0) {
            // 私有项目
            if ($request->input('token')) {
                $storageKey = hash('sha256', $request->input('token'));
                if (!$cacheData = Cache::get($storageKey)) {
                    throw new SecretKeyExpiredException;
                }

                if ($project->id != $cacheData['project_id']) {
                    // 秘钥对应的项目id应该和请求的项目id一致
                    Cache::forget($storageKey);

                    throw ValidationException::withMessages([
                        'project_id' => '请求失败，您传递的信息有误。',
                    ]);
                }

                // 更新缓存时间
                Cache::put($storageKey, $cacheData, 7200);
            } elseif (Auth::guard('api')->check()) {
                if (!ProjectMemberRepository::inThisProject($project->id, Auth::guard('api')->id())) {
                    throw ValidationException::withMessages([
                        'project_id' => '请求失败，您传递的信息有误。'
                    ]);
                }
            } else {
                throw new SecretKeyExpiredException;
            }
        }

        return $project;
    }
}
