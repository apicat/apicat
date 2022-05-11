<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\Project\ApiCommonUrlRepository;

class ApiUrlController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api', 'in.this.project']);
    }

    public function list(Request $request)
    {
        return response()->json([
            'status' => 0,
            'msg' => '',
            'data' => ApiCommonUrlRepository::list($request->input('project_id'))
        ]);
    }

    public function remove(Request $request)
    {
        $request->validate([
            'url_id' => ['required', 'integer', 'min:1'],
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '无法删除',
            ]);
        }

        if (!ApiCommonUrlRepository::remove($request->input('project_id'), $request->input('url_id'))) {
            throw ValidationException::withMessages([
                'url_id' => '删除失败',
            ]);
        }

        return response()->json(['status' => 0, 'msg' => '删除成功']);
    }
}
