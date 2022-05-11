<?php

namespace App\Http\Middleware;

use App\Repositories\Project\ProjectMemberRepository;
use App\Repositories\Project\ProjectRepository;
use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

class InThisProject
{
    /**
     * Handle an incoming request.
     *
     * @param Request $request
     * @param Closure(\Illuminate\Http\Request): (\Illuminate\Http\Response|\Illuminate\Http\RedirectResponse)  $next
     * @return mixed
     */
    public function handle(Request $request, Closure $next)
    {
        if (!$request->input('project_id') or !is_numeric($request->input('project_id'))) {
            return response()->json(['status' => -1, 'msg' => '请求失败，您传递的信息有误。']);
        }

        $project = ProjectRepository::get($request->input('project_id'));
        if (!$project) {
            return response()->json(['status' => -1, 'msg' => '您访问的项目不存在']);
        }

        if (!ProjectMemberRepository::inThisProject($project->id, Auth::id())) {
            return response()->json(['status' => -1, 'msg' => '您访问的项目不存在']);
        }

        ProjectRepository::active($project);

        return $next($request);
    }
}
