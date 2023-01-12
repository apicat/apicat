<?php

namespace App\Http\Controllers\Api;

use App\Exceptions\NotFoundException;
use App\Exceptions\SecretKeyExpiredException;
use App\Http\Controllers\Controller;
use App\Modules\EditorJsonToHtml\Parser;
use App\Repositories\User\UserRepository;
use App\Repositories\Project\ApiDocRepository;
use App\Repositories\Project\DocShareRepository;
use App\Repositories\Project\ProjectMemberRepository;
use App\Repositories\Project\ProjectRepository;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Str;
use Illuminate\Validation\ValidationException;

class ApiDocNoAuthController extends Controller
{
    /**
     * 文档分享状态
     * @param Request $request
     * @return array
     * @throws NotFoundException
     */
    public function hasShared(Request $request)
    {
        $request->validate([
            'doc_id' => 'required|integer|min:1'
        ]);

        $docShare = DocShareRepository::getByDocId($request->input('doc_id'));

        return [
            'status' => 0,
            'msg' => '',
            'data' => (bool)$docShare
        ];
    }

    /**
     * 文档详情
     * @param Request $request
     * @return array
     * @throws NotFoundException
     * @throws SecretKeyExpiredException
     * @throws ValidationException
     */
    public function detail(Request $request)
    {
        $request->validate([
            'project_id' => 'nullable|integer|min:1',
            'doc_id' => 'required|integer|min:1',
            'token' => 'nullable|string|size:60',
            'format' => 'nullable|string|in:json,html',
            'deleted' => 'nullable|boolean'
        ]);

        $format = $request->input('format') ? $request->input('format') : 'json';
        $deleted = $request->input('deleted') ? $request->input('deleted') : false;

        $doc = $this->getDoc($request, $deleted);

        if ($format == 'json') {
            $content = $doc->content;
        } else {
            $content = $doc->content ? Parser::parse($doc->content, $doc->project_id, $doc->id) : '';
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'id' => $doc->id,
                'title' => $doc->title,
                'content' => $content,
                'updated_time' => $doc->updated_at->format('Y-m-d H:i'),
                'last_updated_by' => UserRepository::name($doc->updated_user_id, true)
            ]
        ];
    }

    /**
     * 文档秘钥校验
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function checkSecretKey(Request $request)
    {
        $request->validate([
            'secret_key' => 'required|string|size:6',
            'doc_id' => 'required|integer|min:1'
        ]);

        $docShare = DocShareRepository::getByDocId($request->input('doc_id'));
        if (!$docShare) {
            throw ValidationException::withMessages([
                'doc_id' => '访问密码不正确'
            ]);
        }

        if ($docShare->secret_key != $request->input('secret_key')) {
            throw ValidationException::withMessages([
                'secret_key' => '访问密码不正确',
            ]);
        }

        $token = Str::random(60);
        $storageKey = hash('sha256', $token);

        if (!Cache::put($storageKey, ['doc_id' => $docShare->doc_id], 7200)) {
            throw ValidationException::withMessages([
                'result' => '验证失败，请稍后重试。',
            ]);
        }

        return ['status' => 0, 'msg' => '', 'data' => $token];
    }

    /**
     * 获取文档
     * 
     * @param Request $request
     * @param boolean $deleted — 是否已删除
     * @return Project
     * @throws SecretKeyExpiredException
     * @throws ValidationException
     */
    protected function getDoc($request, $deleted)
    {
        if (!$doc = ApiDocRepository::getNode($request->input('doc_id'), $deleted)) {
            throw ValidationException::withMessages([
                'doc_id' => '您访问的文档不存在',
            ]);
        }

        // project_id参数优先
        $projectID = $request->input('project_id') ? $request->input('project_id') : $doc->project_id;
        if (!$project = ProjectRepository::get($projectID)) {
            throw ValidationException::withMessages([
                'project_id' => '您访问的项目不存在',
            ]);
        }

        if (Auth::guard('api')->check() and ProjectMemberRepository::inThisProject($project->id, Auth::guard('api')->id())) {
            // 登录状态，且属于此项目
            return $doc;
        }

        if ($project->visibility == 0) {
            // 私有项目
            if (!$request->input('token')) {
                throw ValidationException::withMessages([
                    'doc_id' => '您访问的文档不存在',
                ]);
            }

            $storageKey = hash('sha256', $request->input('token'));
            if (!$cacheData = Cache::get($storageKey)) {
                throw new SecretKeyExpiredException;
            }

            if (isset($cacheData['project_id'])) {
                // 整个项目分享
                if ($project->id != $cacheData['project_id']) {
                    // 秘钥对应的项目id应该和请求的项目id一致
                    Cache::forget($storageKey);

                    throw ValidationException::withMessages([
                        'project_id' => '请求失败，您传递的信息有误。',
                    ]);
                }
            } elseif (isset($cacheData['doc_id'])) {
                // 单篇文档分享
                if ($doc->id != $cacheData['doc_id']) {
                    // 秘钥对应的文档id应该和请求的文档id一致
                    Cache::forget($storageKey);

                    throw ValidationException::withMessages([
                        'doc_id' => '请求失败，您传递的信息有误。',
                    ]);
                }

                $docShare = DocShareRepository::getByDocId($doc->id);
                if (!$docShare) {
                    Cache::forget($storageKey);

                    throw ValidationException::withMessages([
                        'doc_id' => '您访问的文档不存在',
                    ]);
                }
            } else {
                Cache::forget($storageKey);

                throw ValidationException::withMessages([
                    'doc_id' => '请求失败，您传递的信息有误。',
                ]);
            }

            // 更新缓存时间
            Cache::put($storageKey, $cacheData, 7200);
        }

        return $doc;
    }
}
