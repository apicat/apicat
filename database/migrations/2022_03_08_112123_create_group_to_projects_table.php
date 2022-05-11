<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('group_to_projects', function (Blueprint $table) {
            $table->id();
            $table->bigInteger('group_id')->index()->comment('项目分组id');
            $table->bigInteger('project_id')->comment('项目id');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('group_to_projects');
    }
};
