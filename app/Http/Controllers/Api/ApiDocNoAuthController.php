<?php

namespace App\Http\Controllers\Api;

use App\Exceptions\NotFoundException;
use App\Exceptions\SecretKeyExpiredException;
use App\Http\Controllers\Controller;
use App\Modules\EditorJsonToHtml\Parser;
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
    public function doc(Request $request)
    {
        $request->validate([
            'token' => 'required|string|size:60',
            'doc_id' => 'required|integer|min:1'
        ]);

        $docShare = DocShareRepository::getByDocId($request->input('doc_id'));
        if (!$docShare) {
            throw new NotFoundException;
        }

        if (!ProjectRepository::get($docShare->project_id)) {
            throw new NotFoundException;
        }

        $storageKey = hash('sha256', $request->input('token'));
        if (!$cacheData = Cache::get($storageKey)) {
            throw new SecretKeyExpiredException;
        }

        if ($docShare->doc_id != $cacheData['doc_id']) {
            // 秘钥对应的文档id应该和请求的文档id一致
            Cache::forget($storageKey);

            throw ValidationException::withMessages([
                'doc_id' => '请求失败，您传递的信息有误。',
            ]);
        }

        // 更新缓存时间
        Cache::put($storageKey, $cacheData, 7200);

        if (!$doc = ApiDocRepository::getNode($docShare->doc_id)) {
            throw new NotFoundException;
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'id' => $doc->id,
                'title' => $doc->title,
                'content' => $doc->content ? Parser::parse($doc->content, $doc->project_id, $doc->id) : '',
                'updated_time' => $doc->updated_at->format('Y-m-d H:i')
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
}
