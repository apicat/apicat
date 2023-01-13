<?php

namespace App\Http\Controllers\Api;

use App\Exceptions\NotFoundException;
use App\Exceptions\NotLoginException;
use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Str;
use Illuminate\Validation\ValidationException;
use App\Exceptions\SecretKeyExpiredException;
use App\Repositories\Iteration\IterationRepository;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\Project\ProjectMemberRepository;
use App\Repositories\Project\ProjectShareRepository;
use App\Repositories\ApiDoc\ApiDocRepository;
use Illuminate\Support\Facades\Auth;

/**
 * 非必需登录访问的API
 */
class ProjectNoAuthController extends Controller
{
    public function status(Request $request)
    {
        $request->validate([
            'project_id' => 'required|integer|min:1'
        ]);

        $result = [];

        if (!$project = ProjectRepository::get($request->input('project_id'))) {
            throw ValidationException::withMessages([
                'project_id' => '您访问的项目不存在',
            ]);
        }

        if (Auth::guard('api')->check()) {
            $authorityDescArr = ['manage', 'write', 'read'];

            $authority = ProjectMemberRepository::getAuthority($project->id, Auth::guard('api')->id());
            if (is_null($authority)) {
                $result['authority'] = 'none';
            } else {
                $result['authority'] = $authorityDescArr[$authority];
            }
        } else {
            $result['authority'] = 'none';
        }

        $result['visibility'] = $project->visibility ? 'public' : 'private';
        $result['has_shared'] = $project->visibility ? true : ProjectShareRepository::hasShared($project->id);

        return ['status' => 0, 'msg' => '', 'data' => $result];
    }

    public function detail(Request $request)
    {
        $request->validate([
            'token' => 'nullable|string|size:60'
        ]);

        if (!$request->input('project_id') and !$request->input('iteration_id')) {
            throw ValidationException::withMessages([
                'project_id' => '请求失败，您传递的信息有误。',
            ]);
        }

        $project = $this->getProject($request);

        if (Auth::guard('api')->check()) {
            $authorityNameArr = ['管理者', '维护者', '阅读者'];
            $authorityDescArr = ['manage', 'write', 'read'];

            $authorityNum = ProjectMemberRepository::getAuthority($project->id, Auth::guard('api')->id());
            if (is_null($authorityNum)) {
                $authority = 'none';
                $authorityName = '';
            } else {
                $authority = $authorityDescArr[$authorityNum];
                $authorityName = $authorityNameArr[$authorityNum];
            }
        } else {
            $authority = 'none';
            $authorityName = '';
        }

        if ($authority != 'none') {
            if ($share = ProjectShareRepository::getByMemberID($project->id, Auth::guard('api')->id())) {
                $secretKey = $share->secret_key;
            } else {
                $secretKey = '';
            }
        } else {
            $secretKey = '';
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'id' => $project->id,
                'icon' => $project->icon,
                'name' => $project->name,
                'authority' => $authority,
                'authority_name' => $authorityName,
                'visibility' => $project->visibility ? 'public' : 'private',
                'secret_key' => $secretKey,
                'description' => $authority == 'none' ? '' : ($project->description ? $project->description : ''),
                'user_id' => $authority == 'none' ? 0 : $project->user_id
            ]
        ];
    }

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

    public function docTree(Request $request)
    {
        $request->validate([
            'token' => 'nullable|string|size:60'
        ]);

        if ($request->input('iteration_id')) {
            return $this->iterationDocTree($request);
        } elseif ($request->input('project_id')) {
            return $this->projectDocTree($request);
        } else {
            throw ValidationException::withMessages([
                'project_id' => '请求失败，您传递的信息有误。',
            ]);
        }
    }

    public function docSearch(Request $request)
    {
        $request->validate([
            'project_id' => 'required|integer|min:1',
            'token' => 'nullable|string|size:60',
            'keywords' => 'required|string|max:255',
            'iteration_id' => 'nullable|integer|min:1',
        ]);

        $project = $this->getProject($request);

        if (mb_strlen($request->input('keywords')) < 2) {
            return ['status' => 0, 'msg' => '', 'data' => []];
        }

        if ($request->input('iteration_id')) {
            $records = ApiDocRepository::searchNodeFromIteration($request->input('iteration_id'), $project->id, $request->input('keywords'));
        } else {
            $records = ApiDocRepository::searchNode($project->id, $request->input('keywords'));
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => $records
        ];
    }

    /**
     * 获取项目
     * @param Request $request
     * @return Project
     * @throws SecretKeyExpiredException
     * @throws ValidationException
     * @throws NotFoundException
     */
    protected function getProject($request)
    {
        if ($request->input('project_id')) {
            if (!$project = ProjectRepository::get($request->input('project_id'))) {
                throw new NotFoundException;
            }
        } elseif ($request->input('iteration_id')) {
            $projectId = IterationRepository::getProjectIdByIterationId($request->input('iteration_id'));
            if (!$projectId) {
                throw new NotFoundException;
            }

            if (!$project = ProjectRepository::get($projectId)) {
                throw new NotFoundException;
            }
        } else {
            throw new NotFoundException;
        }

        if (Auth::guard('api')->check()) {
            // 登录状态
            if (ProjectMemberRepository::inThisProject($project->id, Auth::guard('api')->id())) {
                // 属于此项目
                return $project;
            }
        }

        if ($project->visibility == 0) {
            // 私有项目
            if (!$request->input('token')) {
                throw ValidationException::withMessages([
                    'project_id' => '您访问的项目不存在',
                ]);
            }

            $storageKey = hash('sha256', $request->input('token'));
            if (!$cacheData = Cache::get($storageKey)) {
                throw new SecretKeyExpiredException;
            }

            if (!isset($cacheData['project_id']) or $project->id != $cacheData['project_id']) {
                // 秘钥对应的项目id应该和请求的项目id一致
                Cache::forget($storageKey);

                throw ValidationException::withMessages([
                    'project_id' => '请求失败，您传递的信息有误。',
                ]);
            }

            // 更新缓存时间
            Cache::put($storageKey, $cacheData, 7200);
        }

        return $project;
    }

    /**
     * 获取项目文档树
     *
     * @param Request $request
     * @return array
     */
    protected function projectDocTree($request)
    {
        $project = $this->getProject($request);

        return [
            'status' => 0,
            'msg' => '',
            'data' => ApiDocRepository::getTree($project->id)
        ];
    }

    /**
     * 获取迭代文档树
     *
     * @param Request $request
     * @return void
     */
    protected function iterationDocTree($request)
    {
        if (!Auth::guard('api')->check()) {
            throw new NotLoginException;
        }

        $iteration = IterationRepository::get($request->input('iteration_id'));
        if (!$iteration) {
            throw ValidationException::withMessages([
                'iteration_id' => '您查看的迭代不存在',
            ]);
        }

        if (!ProjectMemberRepository::inThisProject($iteration->project_id, Auth::guard('api')->id())) {
            throw ValidationException::withMessages([
                'project_id' => '您没有查看迭代API的权限',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => IterationRepository::apiTree($iteration)
        ];
    }
}
