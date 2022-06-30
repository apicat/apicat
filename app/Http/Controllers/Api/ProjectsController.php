<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Modules\Image\Uploader;
use App\Repositories\Project\ApiCommonParamRepository;
use App\Repositories\Project\GroupRepository;
use App\Repositories\Project\ProjectRepository;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Validation\ValidationException;

class ProjectsController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api']);
    }

    /**
     * 项目列表
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function index(Request $request)
    {
        $request->validate([
            'group_id' => 'nullable|integer|min:1'
        ]);

        if ($request->input('group_id')) {
            $group = GroupRepository::find($request->input('group_id'));
            if (!$group or Auth::id() != $group->user_id) {
                throw ValidationException::withMessages([
                    'result' => '分组不存在',
                ]);
            }

            $groupName = $group->name;
            $projects = ProjectRepository::groupList(Auth::id(), $request->input('group_id'));
        } else {
            $groupName = '所有项目';
            $projects = ProjectRepository::list(Auth::id());
        }

        $projectArr = [];
        if ($projects->count() > 0) {
            if (!$request->input('group_id')) {
                $groupIds = GroupRepository::relationship(0, Auth::id());
            }

            $authoritys = ['manage', 'write', 'read'];

            foreach ($projects as $project) {
                $projectArr[] = [
                    'id' => $project->id,
                    'group_id' => $request->input('group_id') ? (int)$request->input('group_id') : ($groupIds[$project->id] ?? 0),
                    'preview_link' => $project->preview_link,
                    'default_link' => $project->default_link,
                    'secret_key' => (isset($project->secret_key) and $project->secret_key) ? $project->secret_key : '',
                    'icon' => $project->icon ?? '',
                    'name' => $project->name,
                    'visibility' => $project->visibility ? 'public' : 'private',
                    'authority' => $authoritys[$project->authority],
                    'authority_name' => $project->authority_name
                ];
            }
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'group_name' => $groupName,
                'projects' => $projectArr
            ]
        ];
    }

    /**
     * 创建项目
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function create(Request $request)
    {
        $request->validate([
            'group_id' => ['nullable', 'integer', 'min:0'],
            'icon_link' => ['nullable', 'string', 'max:255'],
            'name' => ['required', 'string', 'min:2', 'max:255'],
            'visibility' => ['required', 'integer', 'in:0,1'],
            'description' => ['nullable', 'string', 'max:255'],
        ]);

        if ($request->input('group_id')) {
            $group = GroupRepository::find($request->input('group_id'));
            if (!$group or $group->user_id != Auth::id()) {
                throw ValidationException::withMessages([
                    'group_id' => '分组不存在',
                ]);
            }
        }

        if ($project = ProjectRepository::create($request->all())) {
            ApiCommonParamRepository::defaultData($project->id);

            return [
                'status' => 0,
                'msg' => '项目创建成功',
                'data' => [
                    'id' => $project->id,
                    'preview_link' => route('app.index', ['projectID' => $project->id]),
                    'default_link' => route('editor.doc', ['projectID' => $project->id]),
                    'secret_key' => '',
                    'icon' => $project->icon ?? '',
                    'name' => $project->name,
                    'visibility' => (integer)$project->visibility,
                    'authority' => 0,
                    'authority_name' => '管理者'
                ]
            ];
        }

        throw ValidationException::withMessages([
            'result' => '项目创建失败，请稍后重试。',
        ]);
    }

    /**
     * 上传项目图标
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws ValidationException
     */
    public function icon(Request $request)
    {
        $request->validate([
            'icon' => 'required|image|mimes:jpeg,jpg,png|max:1024',
            'cropped_x' => 'required|integer|min:0',
            'cropped_y' => 'required|integer|min:0',
            'cropped_width' => 'required|integer|min:1',
            'cropped_height' => 'required|integer|min:1',
        ]);

        if (!$request->hasFile('icon')) {
            throw ValidationException::withMessages([
                'icon' => '请上传项目图标',
            ]);
        }

        if (!$request->file('icon')->isValid()) {
            throw ValidationException::withMessages([
                'avatar' => '无效的图片',
            ]);
        }

        if ($request->input('cropped_width') != $request->input('cropped_height')) {
            throw ValidationException::withMessages([
                'avatar' => '图片裁减信息有误',
            ]);
        }

        $filename = md5(Auth::id() . '|' . time()) . '.' . $request->file('icon')->extension();

        $uploader = new Uploader($request->file('icon')->path(), $filename, '/images/projects/');

        // 检查图片是否被正确加载
        if (!$uploader->imgMakeResult()) {
            throw ValidationException::withMessages([
                'icon' => '无效的图片',
            ]);
        }

        $res = $uploader->croppedSet(
            $request->input('cropped_width'),
            $request->input('cropped_height'),
            $request->input('cropped_x'),
            $request->input('cropped_y')
        );
        if (!$res) {
            throw ValidationException::withMessages([
                'icon' => '图片裁减信息有误',
            ]);
        }

        $uploader->resizeSet(200, 200);
        if ($url = $uploader->save()) {
            return response()->json([
                'status' => 0,
                'msg' => '图片上传成功',
                'data' => asset($url)
            ]);
        }

        throw ValidationException::withMessages([
            'icon' => '图片上传失败，请稍后重试。',
        ]);
    }
}
