<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Validation\ValidationException;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\Project\ProjectMemberRepository;
use App\Repositories\Project\ProjectShareRepository;

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
                'project_id' => '项目不存在',
            ]);
        }

        if (Auth::guard('api')->check()) {
            $authorityDescription = ['manage', 'write', 'read'];

            $authority = ProjectMemberRepository::getAuthority($project->id, Auth::guard('api')->id());
            if (is_null($authority)) {
                $result['authority'] = 'none';
            } else {
                $result['authority'] = $authorityDescription[$authority];
            }
        } else {
            $result['authority'] = 'none';
        }

        $result['visibility'] = $project->visibility ? 'public' : 'private';
        $result['has_shared'] = $project->visibility ? true : ProjectShareRepository::hasShared($project->id);

        return [
            'status' => 0,
            'msg' => '',
            'data' => $result
        ];
    }
}
