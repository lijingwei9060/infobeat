from infobeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Infobeat normally
        """
        self.render_config_template(
            path=os.path.abspath(self.working_dir) + "/log/*"
        )

        infobeat_proc = self.start_beat()
        self.wait_until(lambda: self.log_contains("infobeat is running"))
        exit_code = infobeat_proc.kill_and_wait()
        assert exit_code == 0
